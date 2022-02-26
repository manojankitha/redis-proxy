
.PHONY: build

install-deps:
	go get -v ./...

build: install-deps
	  go build -o redis-proxy cmd/redis-proxy/main.go

test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down --volumes

run: ## Run container. See README for options
	docker-compose up -d

stop: ## Stop docker container
	docker-compose down -v