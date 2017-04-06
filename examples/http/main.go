package main

import (
	"flag"
	"log"

	"github.com/RussellLuo/gotask"
)

func main() {
	addr := flag.String("addr", ":8080", "the local network address")
	flag.Parse()

	worker := HTTPWorker{
		Registry: map[string]gotask.Constructor{
			"add":   func() gotask.Task { return &Add{} },
			"greet": func() gotask.Task { return &Greet{} },
			"panic": func() gotask.Task { return &Panic{} },
		},
		Addr: *addr,
	}
	err := worker.Start()
	if err != nil {
		log.Fatal(err)
	}
}
