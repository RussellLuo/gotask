package main

import (
	"flag"
	"log"

	"github.com/RussellLuo/gotask"
	"github.com/RussellLuo/gotask/examples/tasks"
)

func main() {
	addr := flag.String("addr", ":8080", "the local network address")
	flag.Parse()

	worker := HTTPWorker{
		Registry: map[string]gotask.Constructor{
			"add":   func() gotask.Task { return &tasks.Add{} },
			"greet": func() gotask.Task { return &tasks.Greet{} },
			"panic": func() gotask.Task { return &tasks.Panic{} },
		},
		Addr: *addr,
	}
	if err := worker.Work(); err != nil {
		log.Fatal(err)
	}
}
