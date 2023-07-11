.PHONY: lint test

lint:
	go vet ./...

test:
	go test -v -race ./...
