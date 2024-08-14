# natsjetstream-poc
PoC of Nats Jetstream

## Commands

```
curl -sf https://binaries.nats.dev/nats-io/natscli/nats@latest | sh

# Publish 100000 messages
./nats bench -s localhost:4222 benchstream --js --pub 1 --msgs=100000

# Manually subscribe to 100000 messages
./nats bench -s localhost:4222 benchstream --js --sub 3 --msgs=100000

# Show list of streams
./nats -s localhost:4222 stream list

# List consumers from subject
./nats con ls benchstream

# Show info for consumer 
./nats con info benchstream telegraf_consumers
```

# Go test

```
# Publish 10 messages
./nats bench -s localhost:4222 benchstream --js --pub 1 --msgs=10

# Consume messages
go run consumer.go

# View consumer stats
./nats con info benchstream processor

```


