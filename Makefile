.PHONY: dc run test lint

dc:
	docker-compose up  --remove-orphans --build

test:
	go test -race ./...

lint:
	golangci-lint run
