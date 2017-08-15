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
	"github.com/go-redis/redis"
)

type Options struct {
	Addr        string
	Queue       string
	Concurrency int
	Connections int
	Interval    time.Duration
}

type RedisWorker struct {
	Registry gotask.Registry
	Opts     Options
}

func (w *RedisWorker) Work() error {
	client := redis.NewClient(&redis.Options{
		Addr:        w.Opts.Addr,
		DB:          0, // use default DB
		PoolSize:    w.Opts.Connections,
		PoolTimeout: 10 * time.Second, // wait 10s for connection
	})

	var wg sync.WaitGroup
	for i := 0; i < w.Opts.Concurrency; i++ {
		wg.Add(1)

		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

			for {
				select {
				case <-sigChan:
					wg.Done()
					return
				default:
					err := w.handle(client)
					if err != nil {
						log.Printf("err: %#v", err)
					}
					time.Sleep(w.Opts.Interval)
				}
			}
		}()
	}

	wg.Wait()
	return nil
}

func (w *RedisWorker) handle(client *redis.Client) error {
	taskSig, err := client.LPop(w.Opts.Queue).Result()
	if err == redis.Nil {
		// Got no task
		return nil
	} else if err != nil {
		return err
	}

	log.Printf("Got task: %s", taskSig)

	sig := gotask.Signature{}
	if err := json.Unmarshal([]byte(taskSig), &sig); err != nil {
		return err
	}

	if err := gotask.Process(w.Registry, &sig); err != nil {
		return err
	}

	return nil
}
