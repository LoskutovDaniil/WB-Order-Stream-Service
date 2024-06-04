package subscribe

import (
	"encoding/json"
	"ex0/model"
	"ex0/cache"
	"ex0/storage"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/nats-io/stan.go"
)

type MySubscribe struct {
	st *storage.Storage
	sc stan.Conn
	cach *cache.Cache
}

func (s *MySubscribe) CloseDB() {
    s.st.Close()
}

func NewSubscribe(st *storage.Storage, sc stan.Conn, cach *cache.Cache) *MySubscribe {
	return &MySubscribe{st: st, sc: sc, cach: cach}
}

func (s *MySubscribe) Start() {
	sub, err := s.sc.Subscribe("my-channel", func(msg *stan.Msg) {
		var m model.Model
		if err := json.Unmarshal(msg.Data, &m); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return
		}
		if err := s.st.AddOrder(m); err != nil {
			log.Printf("Failed to insert order: %v", err)
			return
		}
		
		fmt.Println("Order successfully inserted")

		s.cach.Load(m)
	}, stan.DurableName("my-durable"))
	if err != nil {
		log.Fatalf("Failed to My: %v", err)
	}

	log.Println("Waiting for messages. To exit press CTRL+C")
	sub.IsValid()
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)

	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			sub.Unsubscribe()
			s.sc.Close()
			s.CloseDB()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func (s *MySubscribe) LoadJSONDataIntoDB(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	var order model.Model
	if err := json.NewDecoder(file).Decode(&order); err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	if err := s.st.AddOrder(order); err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}

	fmt.Println("Order successfully inserted from JSON")

	s.cach.Load(order)

	return nil
}