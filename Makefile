.PHONY: build
build:
	go build -o server ./cmd/server.go

.PHONY: run
run: build
	./server

.PHONY: dev
dev:
	@hash reflex 2>/dev/null || go get github.com/cespare/reflex
	@reflex -r '\.go$$' -s go run ./cmd/server.go

.PHONY: test
test:
	@go test

.PHONY: migrate-create
migrate-create:
	@migrate create -ext sql -dir migrations $(name)

.PHONY: migrate-up
migrate-up:
	@migrate -path ./migrations -database $${DSN} up $(n)

.PHONY: migrate-down
migrate-down:
	@migrate -path ./migrations -database $${DSN} down $(n)

.PHONY: fmt
fmt:
	@go fmt

.PHONY: docker.dev-up
docker.dev-up:
	@docker-compose -f docker-compose.dev.yaml up -d

.PHONY: docker.dev-build
docker.dev-build:
	@docker-compose -f docker-compose.dev.yaml build

.PHONY: docker.dev-sh
docker.dev-sh:
	@docker container exec -it comusic_api sh

 
