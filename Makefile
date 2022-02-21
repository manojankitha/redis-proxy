
.PHONY: build

install-deps:
	go get -v ./...

build: install-deps
	  go build -o redis-proxy cmd/redis-proxy/main.go

test:
	go test -v ./...