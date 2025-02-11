package natstream

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type (
	// QueueConfig holds the configuration for a jetstream queue
	QueueConfig struct {
		StreamName string
		Subjects   []string
		Storage    jetstream.StorageType
	}

	// Queue represents a jetstream queue
	Queue struct {
		Stream jetstream.Stream
		Js     jetstream.JetStream
	}

	// ConsumerConfig holds the configuration for a jetstream consumer
	ConsumerConfig struct {
		DurableName   string
		AckPolicy     jetstream.AckPolicy
		MaxDeliver    int
		FilterSubject string
	}
)

// New connets to jetstream with given *nats.Conn and creates a jetstream.Stream with given streamName and subjects
func New(ctx context.Context, nc *nats.Conn, cfg QueueConfig) (*Queue, error) {
	if nc == nil {
		return nil, fmt.Errorf("nats connection is nil")
	}

	if cfg.StreamName == "" {
		return nil, fmt.Errorf("stream name cannot be empty")
	}

	if len(cfg.Subjects) == 0 {
		return nil, fmt.Errorf("at least one subject is required")
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("failed to create jestream client --- %s", err.Error())
	}

	if cfg.Storage.String() == "" {
		cfg.Storage = jetstream.MemoryStorage
	}

	s, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     cfg.StreamName,
		Subjects: cfg.Subjects,
		Storage:  cfg.Storage,
	})

	if err != nil {
		return nil, fmt.Errorf("error creating stream --- %w", err)
	}

	return &Queue{
		Stream: s,
		Js:     js,
	}, nil
}

// RegisterConsumer creates a nats jetstream consumer
func (q Queue) RegisterConsumer(ctx context.Context, cfg ConsumerConfig, handler func(jetstream.Msg)) error {
	if handler == nil {
		return fmt.Errorf("::natstream ---handler function is nil")
	}

	if cfg.DurableName == "" {
		return fmt.Errorf("::natsream --- durable name cannot be empty")
	}

	consumerConfig := jetstream.ConsumerConfig{
		Durable:   cfg.DurableName,
		AckPolicy: cfg.AckPolicy,
	}

	if cfg.MaxDeliver > 0 {
		consumerConfig.MaxDeliver = cfg.MaxDeliver
	}

	if cfg.FilterSubject != "" {
		consumerConfig.FilterSubject = cfg.FilterSubject
	}

	consumer, err := q.Stream.CreateOrUpdateConsumer(ctx, consumerConfig)
	if err != nil {
		return fmt.Errorf("::natsream --- failed to create or update consumer: %w", err)
	}

	log.Printf("::natsream --- consumer %s created", cfg.DurableName)

	cc, err := consumer.Consume(handler)
	if err != nil {
		return fmt.Errorf("::natstream --- failed to start consumer: %w", err)
	}

	log.Printf("::natstream --- consumer %s started", cfg.DurableName)

	go func() {
		<-ctx.Done()
		cc.Stop()
	}()
	return nil
}

// Close closes the jetstream connection
func (q *Queue) Close() error {
	q.Js.Conn().Close()
	return nil
}
