# nats-map

Listens to either a NATS or STAN channel and runs a configurable executor for
each message.

## Testing

Run a nats-streaming-server instance in one terminal.

``` shell
$ go get github.com/nats-io/nats-streaming-server
$ nats-streaming-server
[265074] 2020/01/07 11:09:38.675746 [INF] STREAM: Starting nats-streaming-server[test-cluster] version 0.16.2
[265074] 2020/01/07 11:09:38.676353 [INF] STREAM: ServerID: ubReUv2Nky2YFjlbzvWudg
[265074] 2020/01/07 11:09:38.676361 [INF] STREAM: Go version: go1.13.5
[265074] 2020/01/07 11:09:38.676365 [INF] STREAM: Git commit: [not set]
[265074] 2020/01/07 11:09:38.678069 [INF] Starting nats-server version 2.0.4
[265074] 2020/01/07 11:09:38.678133 [INF] Git commit [not set]
[265074] 2020/01/07 11:09:38.678587 [INF] Listening for client connections on 0.0.0.0:4222
[265074] 2020/01/07 11:09:38.678624 [INF] Server id is NDKALSGCMVO3FCMP3CBQ56XHCJRJC7UKX2U426GAWYLWH4A6DVIUCE5W
[265074] 2020/01/07 11:09:38.678631 [INF] Server is ready
[265074] 2020/01/07 11:09:38.708599 [INF] STREAM: Recovering the state...
[265074] 2020/01/07 11:09:38.708707 [INF] STREAM: No recovered state
[265074] 2020/01/07 11:09:38.961856 [INF] STREAM: Message store is MEMORY
[265074] 2020/01/07 11:09:38.962029 [INF] STREAM: ---------- Store Limits ----------
[265074] 2020/01/07 11:09:38.962051 [INF] STREAM: Channels:                  100 *
[265074] 2020/01/07 11:09:38.962060 [INF] STREAM: --------- Channels Limits --------
[265074] 2020/01/07 11:09:38.962070 [INF] STREAM:   Subscriptions:          1000 *
[265074] 2020/01/07 11:09:38.962078 [INF] STREAM:   Messages     :       1000000 *
[265074] 2020/01/07 11:09:38.962086 [INF] STREAM:   Bytes        :     976.56 MB *
[265074] 2020/01/07 11:09:38.962099 [INF] STREAM:   Age          :     unlimited *
[265074] 2020/01/07 11:09:38.962107 [INF] STREAM:   Inactivity   :     unlimited *
[265074] 2020/01/07 11:09:38.962115 [INF] STREAM: ----------------------------------
[265074] 2020/01/07 11:09:38.962123 [INF] STREAM: Streaming Server is ready
```

From the current `nats-map` directory, run a subscriber in one terminal.

``` shell
$ ./scripts/sub.sh
Listening on [nats-map-test], clientID=[nats-map-test-sub-1], qgroup=[] durable=[test]
```

Run the publisher from another terminal.

The publisher will flip a coin between sending an event that has an
"ObjectCreated" or an "ObjectAccessed" event.

The preconfigured mapper will only process the "ObjectCreated" events.

``` shell
$ ./scripts/pub.sh
Published [nats-map-test] : '{"EventName": "s3:ObjectCreated:Put", "Key": "bucket/key"}'
```

The publisher will have additional output.

``` shell
=== {nats-map-test} processing message [seq 1] ===
eventname: "s3:ObjectCreated:Put"
key: "bucket/key"
missing: "null"
```

Distribution Statement "A" (Approved for Public Release, Distribution
Unlimited).
