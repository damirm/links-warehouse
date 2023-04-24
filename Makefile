OUT_DIR := ./out
MIGRATIONS_DIR := $(shell pwd)/internal/postgres/migrations
DOCKER_COMPOSE_FILE := $(shell pwd)/deployments/docker-compose/compose.yaml

.PHONY: test clean sqlc compile recompile docker-build migration

sqlc:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

test: sqlc
	go test -v -race -cover ./...

gh-action-tests:
	act -v

clean:
	rm -rf $(OUT_DIR)/*

compile:
	./scripts/build.sh ./cmd/warehouse ./out/links-warehouse

recompile: sqlc compile

docker-build:
	./scripts/build-docker-image.sh build/docker/Dockerfile links-warehouse "$(shell scripts/git-version.sh)" build

migration:
	@read -p "Enter migration name: " name; \
	docker run --rm -v $(MIGRATIONS_DIR):/migrations migrate/migrate create -ext sql -dir /migrations $$name

compose-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up
compose-build-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up --build
compose-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
compose-psql:
	docker-compose -f $(DOCKER_COMPOSE_FILE) exec postgres psql -U user -d test
