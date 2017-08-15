package main

import (
	"flag"
	"log"
	"strings"

	"github.com/RussellLuo/gotask/examples/tasks"
)

type StringArray []string

func (a *StringArray) Set(s string) error {
	*a = append(*a, s)
	return nil
}

func (a *StringArray) String() string {
	return strings.Join(*a, ",")
}

func main() {
	topic := flag.String("topic", "", "NSQ topic")
	channel := flag.String("channel", "", "NSQ channel")
	maxInFlight := flag.Int("max-in-flight", 200, "max number of messages to allow in flight")
	concurrency := flag.Int("concurrency", 1, "the number of goroutines to spawn for message handling")
	nsqdTCPAddrs := StringArray{}
	lookupdHTTPAddrs := StringArray{}
	flag.Var(&nsqdTCPAddrs, "nsqd-tcp-address", "nsqd TCP address (may be given multiple times)")
	flag.Var(&lookupdHTTPAddrs, "lookupd-http-address", "lookupd HTTP address (may be given multiple times)")

	flag.Parse()

	if *topic == "" {
		log.Fatal("--topic is required")
	}

	if *channel == "" {
		*channel = *topic
	}

	if len(nsqdTCPAddrs) == 0 && len(lookupdHTTPAddrs) == 0 {
		log.Fatal("--nsqd-tcp-address or --lookupd-http-address required")
	}
	if len(nsqdTCPAddrs) > 0 && len(lookupdHTTPAddrs) > 0 {
		log.Fatal("use --nsqd-tcp-address or --lookupd-http-address not both")
	}

	worker := NSQWorker{
		Registry: tasks.Registry,
		Opts: Options{
			Topic:            *topic,
			Channel:          *channel,
			MaxInFlight:      *maxInFlight,
			Concurrency:      *concurrency,
			NSQDTCPAddrs:     nsqdTCPAddrs,
			LookupdHTTPAddrs: lookupdHTTPAddrs,
		},
	}
	if err := worker.Work(); err != nil {
		log.Fatal(err)
	}
}
