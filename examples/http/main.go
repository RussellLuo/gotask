package main

import (
	"flag"
	"log"

	"github.com/RussellLuo/gotask/examples/tasks"
)

func main() {
	addr := flag.String("addr", ":8080", "the local network address")
	flag.Parse()

	worker := HTTPWorker{
		Registry: tasks.Registry,
		Addr:     *addr,
	}
	if err := worker.Work(); err != nil {
		log.Fatal(err)
	}
}
