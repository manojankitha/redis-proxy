

.PHONY: build

install-deps:
	go get -v ./...

build: install-deps
	cd cmd/redis-proxy && \
	  go build

test:
	go test -v ./...