A nats jetstream golang publisher/consumer utility

## getting started

```shell
go get github.com/briannbig/natstream
```

### Example usage

```go
import (
    
	"context"
	"os"
    "time"
    
    "github.com/nats-io/nats.go"
    "github.com/briannbig/natstream"
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
```