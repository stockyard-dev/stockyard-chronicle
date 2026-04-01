build:
	CGO_ENABLED=0 go build -o chronicle ./cmd/chronicle/

run: build
	./chronicle

test:
	go test ./...

clean:
	rm -f chronicle

.PHONY: build run test clean
