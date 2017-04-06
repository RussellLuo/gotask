package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RussellLuo/gotask"
	"github.com/nsqio/go-nsq"
)

type Options struct {
	Topic            string
	Channel          string
	MaxInFlight      int
	Concurrency      int
	NSQDTCPAddrs     []string
	LookupdHTTPAddrs []string
}

type NSQWorker struct {
	Registry map[string]gotask.Constructor
	Opts     Options
}

func (w *NSQWorker) Start() error {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = w.Opts.MaxInFlight

	consumer, err := nsq.NewConsumer(w.Opts.Topic, w.Opts.Channel, cfg)
	if err != nil {
		return err
	}

	consumer.AddConcurrentHandlers(w, w.Opts.Concurrency)

	err = consumer.ConnectToNSQDs(w.Opts.NSQDTCPAddrs)
	if err != nil {
		return err
	}

	err = consumer.ConnectToNSQLookupds(w.Opts.LookupdHTTPAddrs)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-consumer.StopChan:
			return nil
		case <-sigChan:
			consumer.Stop()
		}
	}
}

func (w *NSQWorker) HandleMessage(m *nsq.Message) error {
	log.Printf("Received message: %s", m.Body)

	sig := gotask.Signature{}
	err := json.Unmarshal(m.Body, &sig)
	if err != nil {
		log.Printf("err: %#v", err)
		return err
	}

	err = gotask.Process(w.Registry, &sig)
	if err != nil {
		log.Printf("err: %#v", err)
	}

	return err
}
