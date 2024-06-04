package main

import (
	"ex0/cache"
	"ex0/server"
	"ex0/storage"
	"ex0/subscribe"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
)

func main() {
	clusterID := "test-cluster"
	clientID := "sub-client"
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer sc.Close()

	db, err := storage.NewPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	st := storage.NewStorage(db)

	cache := cache.NewCache(st)
	cache.Init()

	subscriber := subscribe.NewSubscribe(st, sc, cache)

	server := &server.Server{Cache: cache}

	go subscriber.Start()

	go server.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
}