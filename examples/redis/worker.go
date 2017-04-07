package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RussellLuo/gotask"
	"github.com/garyburd/redigo/redis"
)

type Options struct {
	Addr        string
	Queue       string
	Concurrency int
	Interval    time.Duration
}

type RedisWorker struct {
	Registry map[string]gotask.Constructor
	Opts     Options
}

func (w *RedisWorker) Start() error {
	conns := []redis.Conn{}
	for i := 0; i < w.Opts.Concurrency; i++ {
		conn, err := redis.Dial("tcp", w.Opts.Addr)
		if err != nil {
			return err
		}
		defer conn.Close()
		conns = append(conns, conn)
	}

	stopChan := make(chan os.Signal, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for _, conn := range conns {
		go func() {
			for {
				select {
				case <-sigChan:
					stopChan <- syscall.SIGTERM
					return
				default:
					err := w.handle(conn)
					if err != nil {
						log.Printf("err: %#v", err)
					}
					time.Sleep(w.Opts.Interval)
				}
			}
		}()
	}

	<-stopChan
	return nil
}

func (w *RedisWorker) handle(conn redis.Conn) error {
	bytes, err := conn.Do("LPOP", w.Opts.Queue)
	if err != nil {
		return err
	}

	// Got no task
	if bytes == nil {
		return nil
	}
	log.Printf("Got task: %s", bytes)

	sig := gotask.Signature{}
	err = json.Unmarshal(bytes.([]byte), &sig)
	if err != nil {
		return err
	}

	err = gotask.Process(w.Registry, &sig)
	if err != nil {
		return err
	}

	return nil
}
