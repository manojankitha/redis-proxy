
.PHONY: build

install-deps:
	go get -v ./...

build: install-deps
	  go build -o redis-proxy cmd/redis-proxy/main.go

test:
	go test -v ./...

run: ## Run container. See README for options
	docker-compose up -d

stop: ## Stop docker container
	docker-compose down -v