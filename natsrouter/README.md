# natsrouter

Provide a simple webhook server that listens for messages and passes them on to
a NATS queue.

The server takes a "Key" parameter if present and uses that to route messages.

## Metrics

In a future effort, the prometheus middleware will be enabled for more precise
profiling of the webhook server.

At a quick glance,

```shell
$ nomad logs -job service/webhook | jq -R 'fromjson? | select(type == "object") | .latency_human'
"269.228µs"
"221.544µs"
"8.876µs"
"24.644µs"
"24.644µs"
"277.595µs"
"10.204µs"
"256.148µs"
"16.377µs"
"230.971µs"
"10.207µs"
"317.493µs"
"244.577µs"
"227.865µs"
"43.823µs"
"4.836µs"
"12.263µs"
"258.207µs"
"319.458µs"
"204.524µs"
"181.707µs"
"265.075µs"
"181.221µs"
"39.898µs"
"18.241µs"
"263.255µs"
```

The rough estimate is that for empty messages, latency is 20 nanoseconds whereas
with the additional NATS call, it is closer to 250 nanoseconds.

## Testing

First, build the `natsrouter` utility.

``` shell
$ make
go get github.com/nats-io/nats.go/
go get github.com/spf13/viper
gofmt -s -w main.go
CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w"
```

Next, run the utility against an NATS server.
The default prefix is `test` and is safe to use.

# XXX replace with nats_domain
``` shell
$ NATS_URL=nats://<nats_domain>:4222 ./natsrouter
```

In another terminal, listen to the NATS `test.>` wildcard channel.

# XXX replace with nats_domain
``` shell
$ go get github.com/nats-io/nats.go/examples/nats-sub
$ nats-sub -t -s nats://<nats_domain> "test.>"
```

Hit the `natsrouter` webhook.

```shell
$ curl -H "Content-Type: application/json" -d '{"Key": "data-bysha256/hexdigest"}' http://localhost:9090
```

In the shell with `nats-sub` running, you should see:

``` shell
Listening on [test.>]
2019/12/26 22:20:04 [#1] Received on [test.data-bysha256]: '{"Key":"data-bysha256/hexdigest"}'
```

Distribution Statement "A" (Approved for Public Release, Distribution
Unlimited).
