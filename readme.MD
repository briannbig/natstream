A nats jetstream golang cunsumer utility

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
	q := natstream.New(ctx, nc, "notification", []string{"subject.one","subject.two","subject.three"})

	svc := NotificationService()

    // Create consumer
	q.RegisterConsumer(ctx, svc.Notify, consumerDurableName)

    //... rest of your code

}
```