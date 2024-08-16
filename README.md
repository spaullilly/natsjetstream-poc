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

- With producer
  - Produce **10** messages on stream **ptest** subject **poc1**
  
Command:
```
go run main.go -stream=ptest -producer -subject=poc1 -msgcount=100
```

Output:
```
creating 100 messages
```
  - Check stream

Command:
```
./nats -s localhost:4222 stream list
```

Output:

```
╭─────────────────────────────────────────────────────────────────────────────────────╮
│                                       Streams                                       │
├─────────────┬─────────────┬─────────────────────┬──────────┬─────────┬──────────────┤
│ Name        │ Description │ Created             │ Messages │ Size    │ Last Message │
├─────────────┼─────────────┼─────────────────────┼──────────┼─────────┼──────────────┤
│ ptest       │             │ 2024-08-15 15:40:27 │ 110      │ 5.2 KiB │ 4m45s        │
│ benchstream │             │ 2024-08-14 15:03:11 │ 165      │ 22 KiB  │ 9m36s        │
╰─────────────┴─────────────┴─────────────────────┴──────────┴─────────┴──────────────╯
```

  - Consume
  
Command:
```
 go run main.go -fetch=45 -stream=ptest -ack -acktype=double -consumer=processor
```

Ouput:

```
subject: poc1, message: test message 36
subject: poc1, message: test message 37
subject: poc1, message: test message 38
subject: poc1, message: test message 39
subject: poc1, message: test message 40
subject: poc1, message: test message 41
subject: poc1, message: test message 42
subject: poc1, message: test message 43
subject: poc1, message: test message 44
subject: poc1, message: test message 45
subject: poc1, message: test message 46
subject: poc1, message: test message 47
subject: poc1, message: test message 48
subject: poc1, message: test message 49
subject: poc1, message: test message 50
subject: poc1, message: test message 51
subject: poc1, message: test message 52
subject: poc1, message: test message 53
subject: poc1, message: test message 54
subject: poc1, message: test message 55
subject: poc1, message: test message 56
subject: poc1, message: test message 57
subject: poc1, message: test message 58
subject: poc1, message: test message 59
subject: poc1, message: test message 60
subject: poc1, message: test message 61
subject: poc1, message: test message 62
subject: poc1, message: test message 63
subject: poc1, message: test message 64
subject: poc1, message: test message 65
subject: poc1, message: test message 66
subject: poc1, message: test message 67
subject: poc1, message: test message 68
subject: poc1, message: test message 69
subject: poc1, message: test message 70
subject: poc1, message: test message 71
subject: poc1, message: test message 72
subject: poc1, message: test message 73
subject: poc1, message: test message 74
subject: poc1, message: test message 75
subject: poc1, message: test message 76
subject: poc1, message: test message 77
subject: poc1, message: test message 78
subject: poc1, message: test message 79
subject: poc1, message: test message 80
total received: 45
```

## Retention
- RetentionPolicy: https://docs.nats.io/nats-concepts/jetstream/streams#retentionpolicy
- DiscardPolicy: https://docs.nats.io/nats-concepts/jetstream/streams#discardpolicy

In the stream config, I set MaxAge to 1min. The old messages are deleted after that time. This is default behavior.

```
> go run main.go -stream=ptest -producer -subject=poc1 -msgcount=100 -msgage=60
> ./nats -s localhost:4222 stream info ptest
Information for Stream ptest created 2024-08-15 15:40:27

              Subjects: poc1
              Replicas: 1
               Storage: File

Options:

             Retention: Limits
       Acknowledgments: true
        Discard Policy: Old
      Duplicate Window: 2m0s
     Allows Msg Delete: true
          Allows Purge: true
        Allows Rollups: false

Limits:

      Maximum Messages: unlimited
   Maximum Per Subject: unlimited
         Maximum Bytes: unlimited
           Maximum Age: 1m0s
  Maximum Message Size: unlimited
     Maximum Consumers: unlimited

State:

              Messages: 110
                 Bytes: 5.2 KiB
        First Sequence: 1 @ 2024-08-15 15:44:24
         Last Sequence: 110 @ 2024-08-15 15:47:57
      Active Consumers: 1
    Number of Subjects: 1
```

# NATs non-jetstream test

```
nats bench foo --pub 1 --sub 1 --size 16
```

# Benchmark
https://docs.nats.io/using-nats/nats-tools/nats_cli/natsbench

- Sending 100000 messages
```
./nats bench -s localhost:4222 benchstreamtest --js --pub 1 --msgs=100000
15:54:59 Starting JetStream benchmark [subject=benchstreamtest, multisubject=false, multisubjectmax=100000, js=true, msgs=100,000, msgsize=128 B, pubs=1, subs=0, stream=benchstream, maxbytes=1.0 GiB, storage=file, syncpub=false, pubbatch=100, jstimeout=30s, pull=false, consumerbatch=100, push=false, consumername=natscli-bench, replicas=1, purge=false, pubsleep=0s, subsleep=0s, dedup=false, dedupwindow=2m0s]
15:54:59 Starting publisher, publishing 100,000 messages
Finished      1s [==========================================================================================================================================================================================================] 100%

Pub stats: 53,347 msgs/sec ~ 6.51 MB/sec
```


# Interesting errors
- When trying to use an existing subject from one stream on a new stream:

  - A Stream can have multiple subjects
  - A Consumer can filter on subjects or not
  - A Stream can have multiple Consumers
  - A Subject CAN NOT be in 2 Streams, unless each Stream is on a different Account

```
error creating or updating stream nats: API error: code=400 err_code=10065 description=subjects overlap with an existing stream
error producing messages nats: API error: code=400 err_code=10065 description=subjects overlap with an existing stream
```

- When trying to publish:
  - I think this is an issue because `nats bench` is trying to create the stream and it already existed. Unlike my code, that uses a function to CreateOrUpdate.
```
./nats bench -s localhost:4222 benchstreamtest --js --pub 1 --msgs=100000
15:54:00 Starting JetStream benchmark [subject=benchstreamtest, multisubject=false, multisubjectmax=100000, js=true, msgs=100,000, msgsize=128 B, pubs=1, subs=0, stream=benchstream, maxbytes=1.0 GiB, storage=file, syncpub=false, pubbatch=100, jstimeout=30s, pull=false, consumerbatch=100, push=false, consumername=natscli-bench, replicas=1, purge=false, pubsleep=0s, subsleep=0s, dedup=false, dedupwindow=2m0s]
15:54:00 nats: stream name already in use. If you want to delete and re-define the stream use `nats stream delete benchstream`.
```