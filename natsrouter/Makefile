all: main.go
	go get github.com/labstack/echo/...
	go get github.com/nats-io/nats.go/
	go get github.com/spf13/viper

	gofmt -s -w main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w"
