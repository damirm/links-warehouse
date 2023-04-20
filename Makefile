OUT_DIR := ./out
MIGRATIONS_DIR := $(shell pwd)/internal/postgres/migrations

.PHONY: test clean sqlc compile docker-build migration

sqlc:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

test: sqlc
	go test -v -race -cover ./...

gh-action-tests:
	act -v

clean:
	rm -rf $(OUT_DIR)/*

compile: sqlc
	./scripts/build.sh ./cmd/warehouse ./out/links-warehouse

docker-build:
	./scripts/build-docker-image.sh build/docker/Dockerfile links-warehouse "$(shell scripts/git-version.sh)" build

migration:
	@read -p "Enter migration name: " name; \
	docker run --rm -v $(MIGRATIONS_DIR):/migrations migrate/migrate create -ext sql -dir /migrations $$name
