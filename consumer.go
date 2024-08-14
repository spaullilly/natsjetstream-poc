package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	fetchCntPtr := flag.Int("fetch", 5, "number of messages to fetch")
	//ackPrt := flag.Bool("ack", true, "acknowledge messages")
	flag.Parse()

	// Use the env variable if running in the container, otherwise use the default.
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	// Create an unauthenticated connection to NATS.
	nc, _ := nats.Connect(url)
	defer nc.Drain()

	// Access JetStream for managing streams and consumers as well as for
	// publishing and consuming messages to and from the stream.
	js, _ := jetstream.New(nc)

	streamName := "benchstream"

	//	// JetStream API uses context for timeouts and cancellation.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := js.Stream(ctx, streamName)
	if err != nil {
		fmt.Println("stream not found")
		os.Exit(1)
	}

	fmt.Println("stream found")

	// Durable consumers can be created by specifying the Durable name.
	// Durable consumers are not removed automatically regardless of the
	// InactiveThreshold. They can be removed by calling `DeleteConsumer`.
	dur, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable: "processor2",
		//AckPolicy: jetstream.AckAllPolicy,
	})

	// Consume and fetch work the same way for durable consumers.
	received := 0

	msgs, _ := dur.FetchNoWait(*fetchCntPtr)
	for msg := range msgs.Messages() {
		fmt.Printf("received %q message from durable consumer\n", msg.Subject())
		received++
		//if !*ackPrt {
		if err := msg.DoubleAck(ctx); err != nil {
			fmt.Println("error acking message:", err)
		}
		//}
	}

	fmt.Printf("total received: %d\n", received)

	// While ephemeral consumers will be removed after InactiveThreshold, durable
	// consumers have to be removed explicitly if no longer needed.
	//stream.DeleteConsumer(ctx, "processor")

	// Let's try to get the consumer to make sure it's gone.
	//_, err := stream.Consumer(ctx, "processor")

	//fmt.Println("consumer deleted:", errors.Is(err, jetstream.ErrConsumerNotFound))
}
