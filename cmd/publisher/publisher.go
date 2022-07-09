package main

import (
	"flag"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
)

func main() {
	var fileJSON string
	flag.StringVar(&fileJSON, "file", "./resources/model.json", "JSON file")
	flag.Parse()

	file, err := ioutil.ReadFile(fileJSON)
	if err != nil {
		log.Fatalln(err)
	}

	sc, err := stan.Connect("test-cluster", "main")
	if err != nil {
		log.Fatalln(err)
	}
	// Close connection
	defer func(sc stan.Conn) {
		if err = sc.Close(); err != nil {

		}
	}(sc)

	// Simple Synchronous Publisher
	err = sc.Publish("foo", file) // does not return until an ack has been received from NATS Streaming
	if err != nil {
		log.Fatalln(err)
	}
}
