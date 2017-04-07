package main

import (
	"flag"
	"log"
	"time"

	"github.com/RussellLuo/gotask"
	"github.com/RussellLuo/gotask/examples/tasks"
)

func main() {
	addr := flag.String("addr", ":6379", "the TCP address of Redis")
	queue := flag.String("queue", "", "the task queue name in Redis")
	concurrency := flag.Int("concurrency", 1, "the number of goroutines (each has a connection to Redis) to spawn for task handling")
	intervalStr := flag.String("interval", "10ms", "the interval for polling the task queue")

	flag.Parse()

	if *queue == "" {
		log.Fatal("-queue is required")
	}

	interval, err := time.ParseDuration(*intervalStr)
	if err != nil {
		log.Fatal(err)
	}

	worker := RedisWorker{
		Registry: map[string]gotask.Constructor{
			"add":   func() gotask.Task { return &tasks.Add{} },
			"greet": func() gotask.Task { return &tasks.Greet{} },
			"panic": func() gotask.Task { return &tasks.Panic{} },
		},
		Opts: Options{
			Addr:        *addr,
			Queue:       *queue,
			Concurrency: *concurrency,
			Interval:    interval,
		},
	}
	if err := worker.Work(); err != nil {
		log.Fatal(err)
	}
}
