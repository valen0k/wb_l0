package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4"
	"github.com/nats-io/stan.go"
	"github.com/valen0k/wb_l0/internal/app"
	"github.com/valen0k/wb_l0/internal/config"
	"log"
	"net/http"
	"time"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "conf", "./config.json", "service configuration file")
	flag.Parse()

	log.Println("start application")

	log.Println("reading file configuration")
	config, err := config.NewConfig(configFile)
	checkErr(err)

	newApp := app.App{}

	log.Println("connection db")
	newApp.DB, err = config.NewDBConnection()
	checkErr(err)
	defer func(connection *pgx.Conn, ctx context.Context) {
		if err = connection.Close(ctx); err != nil {
			log.Println(err)
		}
	}(newApp.DB, context.Background())

	err = newApp.MemoryRecovery()
	checkErr(err)

	sc, err := stan.Connect("test-cluster", "client-123")
	checkErr(err)
	defer func(sc stan.Conn) {
		if err = sc.Close(); err != nil {
			log.Println(err)
		}
	}(sc)

	// Subscribe starting a specific amount of time in the past (e.g. 30 seconds ago)
	sc.Subscribe("foo", func(msg *stan.Msg) {
		newApp.Rec(msg)
	}, stan.StartAtTimeDelta(30*time.Second))

	server := &http.Server{
		Addr:         config.Server.Host + ":" + config.Server.Port,
		Handler:      newApp.NewHandler(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("app started")
	log.Fatalln(server.ListenAndServe())
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
