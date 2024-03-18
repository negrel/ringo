
.PHONY: test
test:
	go test -v -race ./...

.PHONY: bench
bench:
	go test -v -run=XXX -bench=. -benchmem

.PHONY: lint
lint:
	golangci-lint run --timeout 2m ./...

