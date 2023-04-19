OUT_DIR := ./out

sqlc:
	docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate

test:
	go test -v -race -cover ./...

clean:
	rm -rf $(OUT_DIR)/*

compile:
	./scripts/build.sh ./cmd/warehouse ./out/links-warehouse

docker:
	./scripts/build-docker-image.sh build/docker/Dockerfile links-warehouse "$(shell scripts/git-version.sh)" build
