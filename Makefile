
.PHONY: build

install-deps:
	go get -v ./...

build: install-deps
	  go build -o redis-proxy cmd/redis-proxy/main.go

test: ## assuming its run of mac OS with brew installed.
	brew install golang
	go test -v ./...

run: ## Run container. See README for options
	docker-compose up -d

stop: ## Stop docker container
	docker-compose down -v