.PHONY: all
all: main.go
	go get -u github.com/docopt/docopt-go
	gofmt -s -w main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w"

.PHONY: bundle
bundle: bundle.json
	jq -r '.data.issuing_ca' < bundle.json > bundle.pem
	jq -r '.data.private_key' < bundle.json >> bundle.pem
	jq -r '.data.certificate' < bundle.json >> bundle.pem

.PHONY: clean
clean:
	rm -f bundle.pem
	rm -f service.crt
	rm -f service.key
	rm -f CA.crt
