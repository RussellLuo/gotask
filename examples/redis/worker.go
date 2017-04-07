package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
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

	var wg sync.WaitGroup
	for _, conn := range conns {
		wg.Add(1)

		go func(conn redis.Conn) {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			for {
				select {
				case <-sigChan:
					wg.Done()
					return
				default:
					err := w.handle(conn)
					if err != nil {
						log.Printf("err: %#v", err)
					}
					time.Sleep(w.Opts.Interval)
				}
			}
		}(conn)
	}

	wg.Wait()
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
