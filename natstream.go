package natstream

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type (
	consumer struct {
		con      jetstream.Consumer
		callback func(jetstream.Msg)
	}

	Queue struct {
		Stream jetstream.Stream
		Js     jetstream.JetStream
	}
)

// connets to jetstream with given *nats.Conn and creates a jetstream.Stream with given streamName and subjects
func New(ctx context.Context, nc *nats.Conn, streamName string, subjects []string) Queue {
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Error connecting to jetstream --- %s", err.Error())
	}

	s, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
		Storage:  jetstream.MemoryStorage,
	})
	if err != nil {
		log.Fatalf("error creating stream --- %s", err.Error())
	}

	return Queue{Stream: s, Js: js}
}

// creates a nats jetstream consumer
func (q Queue) RegisterConsumer(ctx context.Context, handler func(jetstream.Msg), consumerDurableName string) error {

	cons, _ := q.Stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   consumerDurableName,
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	return newConsumer(cons, handler)
}

func newConsumer(con jetstream.Consumer, handler func(jetstream.Msg)) error {
	var c = &consumer{
		con:      con,
		callback: handler,
	}
	_, err := c.con.Consume(c.callback)
	return err
}
