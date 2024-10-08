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
> go run main.go -stream=ptest -producer -subject=poc1 -msgcount=100 -msgage=600
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

# Good Docs
- https://docs.nats.io/nats-concepts/jetstream/streams
- https://github.com/ConnectEverything/nats-by-example/tree/main/examples/jetstream/pull-consumer/go


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

# Various Commands
```
go run main.go -fetch=10 -ack -acktype=double -stream=benchstream -consumer=processor2

go run main.go -fetch=10 -stream=benchstream -consumer=processor2


# Produce and consume via code
go run main.go -stream=ptest -producer -subject=poc1 -msgcount=10
go run main.go -fetch=5 -stream=ptest -consumer=processor

```


# Subjects

## Montoring Data

- Logs: o11y.rdu1.logs.t0-gg1-a1-01-r001-rdu1
- Metrics: o11y.rdu1.metrics.cw-cumulus-exporter.t0-gg1-a1-01-r001-rdu1

### Produce

- Subject: `o11y.rdu1.logs.t0-gg1-a1-01-r001-rdu1`
```
go run main.go -stream=subtest -producer -subject=o11y.rdu1.logs.t0-gg1-a1-01-r001-rdu1 -msgcount=100 -msgage=600
```

- Subject: `o11y.rdu1.metrics.cw-cumulus-exporter.t0-gg1-a1-01-r001-rdu1`
```
go run main.go -stream=subtest -producer -subject=o11y.rdu1.metrics.cw-cumulus-exporter.t0-gg1-a1-01-r001-rdu1 -msgcount=100 -msgage=600
```

### Consume
- Consume from first
```
go run main.go -fetch=5 -stream=subtest -ack -acktype=double -consumer=sub1 -subjects=o11y.rdu1.logs.t0-gg1-a1-01-r001-rdu1
```

- Consume from second
```
go run main.go -fetch=5 -stream=subtest -ack -acktype=double -consumer=sub1 -subjects=o11y.rdu1.metrics.cw-cumulus-exporter.t0-gg1-a1-01-r001-rdu1
```

- Consume from both
```
go run main.go -fetch=200 -stream=subtest -ack -acktype=double -consumer=sub1 -subjects=o11y.rdu1.logs.t0-gg1-a1-01-r001-rdu1,o11y.rdu1.metrics.cw-cumulus-exporter.t0-gg1-a1-01-r001-rdu1
```

- New consumer from both
```
go run main.go -fetch=200 -stream=subtest -ack -acktype=double -consumer=sub2 -subjects=o11y.rdu1.logs.t0-gg1-a1-01-r001-rdu1,o11y.rdu1.metrics.cw-cumulus-exporter.t0-gg1-a1-01-r001-rdu1
```

- New consumer from both with wildcard
```
go run main.go -fetch=200 -stream=subtest -ack -acktype=double -consumer=sub3 -subjects='o11y.rdu1.>'
```

### Observation
1. After producing to both, logs first then metrics:
  - Using the same consumer, it seems work sequentially. If I consume from the metrics first, then the first logs may or may not be consumable.
  - If I consume from both subjects it works
1. Consuming from both with a new consumer works as expected.

#### Conclusion
THe test proves that with our long term running services we are able to utilize subject-based messaging well. Both listing out the subjects and using the wildcard worked perfectly.


# Queue Groups

Install telegraf for this: https://docs.influxdata.com/telegraf/v1/install/


Test: 
- Spin up 3 telegraf instances all in the same queue group

- Produce messages

```
go run main.go -stream=queuetest -producer -subject=queue.o11y.rdu1.logs.t0-gg1-a1-01-r001-rdu1 -msgcount=100 -msgage=1200
```

- Observe how many each received
  - All three received unique messages
- Turn one off and produce messages:
  - Both remaining got messages (new only)
- Turn back on the 3rd, produce and observe:
  - All three received unique messages (new only)
- Turn them all off, produce message with new string, wait 5 min, produce message with new string, turn them on:
  - For this I didn't need to wait. I turned them all off, produced 30, turned one on and it replayed the previous as well

#### Conclusion
If we didn't care about disaster recovery, we could use telegraf. Issue is, telegraf deletes its consumer when it stops. We could open a ticket to get durable consumers added but the purpose of telegraf is to be ephemeral. If we had a resilient setup (multiple hosts) we could scale up telegraf and use as is. 1 consumer would be online at all times in that queue group.