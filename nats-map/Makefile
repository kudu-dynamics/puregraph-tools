all: main.go
	go mod tidy
	gofmt -s -w main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w"
