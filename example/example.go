package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/briannbig/natstream"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	consumerDurableName = "notifications.consumer"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}
	nc, err := nats.Connect(url)
	defer nc.Close()

	if err != nil {
		log.Printf("could not connect to queue --- %s", err.Error())
		os.Exit(1)
	}

	// initialize jetstream
	q, err := natstream.New(ctx, nc, natstream.QueueConfig{
		StreamName: "notification-stream",
		Subjects:   []string{"subject.one", "subject.two", "subject.three"},
		Storage:    jetstream.MemoryStorage,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer q.Close()

	svc := NotificationService()

	// Create consumer
	q.RegisterConsumer(ctx, natstream.ConsumerConfig{
		DurableName: consumerDurableName,
		AckPolicy:   jetstream.AckAllPolicy,
	}, svc.Notify)

	//... rest of your code

}

type (
	service interface {
		Notify(jetstream.Msg)
	}
	notification struct{}
)

// Notify implements service.
func (n notification) Notify(msg jetstream.Msg) {
	log.Print(msg)
}

func NotificationService() service {
	return notification{}
}
