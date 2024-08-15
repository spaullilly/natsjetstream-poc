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

type ConsumerConfig struct {
	Name       string
	FetchCount int
	Ack        bool
	AckType    string
}

type ProducerConfig struct {
	MsgNum  int
	Subject string
}

type Operator struct {
	Ctx            context.Context
	CtxCancel      context.CancelFunc
	NatsURL        string
	JetStream      jetstream.JetStream
	StreamConfig   jetstream.StreamConfig
	ConsumerConfig ConsumerConfig
	ProducerConfig ProducerConfig
}

func NewOperator() *Operator {
	flag.Parse()
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	o := Operator{
		Ctx:       ctx,
		CtxCancel: cancel,
		NatsURL:   url,
	}

	return &o
}

func main() {
	fetchCntPtr := flag.Int("fetch", 5, "number of messages to fetch")
	ackPrt := flag.Bool("ack", false, "acknowledge messages")
	ackType := flag.String("acktype", "single", "acknowledge type (single, double)")
	streamPtr := flag.String("stream", "benchstream", "stream name")
	consumerPtr := flag.String("consumer", "", "consumer name")
	producerPtr := flag.Bool("producer", false, "producer name")
	producerSubjectPtr := flag.String("subject", "poc", "subject name")
	producerMsgCountPtr := flag.Int("msgcount", 10, "number of messages to produce")
	producerMsgAgePtr := flag.Int("msgage", 60, "message age in seconds")
	flag.Parse()

	o := NewOperator()
	o.StreamConfig.Name = *streamPtr
	o.StreamConfig.Subjects = []string{*producerSubjectPtr}
	o.StreamConfig.MaxAge = time.Duration(*producerMsgAgePtr) * time.Second

	// Create an unauthenticated connection to NATS.
	nc, _ := nats.Connect(o.NatsURL)
	defer nc.Drain()

	// Access JetStream for managing streams and consumers as well as for
	// publishing and consuming messages to and from the stream.
	js, err := jetstream.New(nc)
	if err != nil {
		fmt.Println("error creating jetstream", err)
		os.Exit(1)
	}
	o.JetStream = js

	//	// JetStream API uses context for timeouts and cancellation.
	o.Ctx, o.CtxCancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer o.CtxCancel()

	if *consumerPtr != "" && *producerPtr {
		fmt.Println("consumer and producer cannot be set at the same time")
		os.Exit(1)
	} else if *consumerPtr != "" {
		o.ConsumerConfig = ConsumerConfig{
			Name:       *consumerPtr,
			FetchCount: *fetchCntPtr,
			Ack:        *ackPrt,
			AckType:    *ackType,
		}

		if err := o.Consume(); err != nil {
			fmt.Println("error consuming messages", err)
		}
	} else if *producerPtr {
		o.ProducerConfig = ProducerConfig{
			MsgNum:  *producerMsgCountPtr,
			Subject: *producerSubjectPtr,
		}
		if err := o.Produce(); err != nil {
			fmt.Println("error producing messages", err)
		}
	}

	// While ephemeral consumers will be removed after InactiveThreshold, durable
	// consumers have to be removed explicitly if no longer needed.
	//stream.DeleteConsumer(ctx, "processor")

	// Let's try to get the consumer to make sure it's gone.
	//_, err := stream.Consumer(ctx, "processor")

	//fmt.Println("consumer deleted:", errors.Is(err, jetstream.ErrConsumerNotFound))
}

func (o *Operator) Consume() error {
	stream, err := o.JetStream.Stream(o.Ctx, o.StreamConfig.Name)
	if err != nil {
		fmt.Printf("%s stream not found\n", o.StreamConfig.Name)
		os.Exit(1)
	}

	dur, _ := stream.CreateOrUpdateConsumer(o.Ctx, jetstream.ConsumerConfig{
		Durable: o.ConsumerConfig.Name,
		//AckPolicy: jetstream.AckAllPolicy,
	})

	// Consume and fetch work the same way for durable consumers.
	received := 0

	msgs, _ := dur.FetchNoWait(o.ConsumerConfig.FetchCount)
	for msg := range msgs.Messages() {
		fmt.Printf("subject: %s, message: %s\n", msg.Subject(), string(msg.Data()))
		received++
		if o.ConsumerConfig.Ack {
			if o.ConsumerConfig.AckType == "single" {
				if err := msg.Ack(); err != nil {
					fmt.Println("error single acking message:", err)
				}
			} else {
				if err := msg.DoubleAck(o.Ctx); err != nil {
					fmt.Println("error double acking message:", err)
				}
			}
		}
	}

	fmt.Printf("total received: %d\n", received)

	return nil
}

func (o *Operator) Produce() error {
	_, err := o.JetStream.CreateOrUpdateStream(o.Ctx, o.StreamConfig)
	if err != nil {
		fmt.Println("error creating or updating stream", err)
		return err
	}
	fmt.Printf("creating %d messages\n", o.ProducerConfig.MsgNum)
	for i := 1; i <= o.ProducerConfig.MsgNum; i++ {
		_, err = o.JetStream.Publish(o.Ctx, o.ProducerConfig.Subject, []byte(fmt.Sprintf("test message %d", i)))
		if err != nil {
			fmt.Println("error publishing message", err)
			return err
		}
	}

	return nil
}
