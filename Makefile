BINARY := 10-nono-nri
CMD := ./cmd/nono-nri

.PHONY: build test clean docker-build docker-load-kind

build:
	go build -o $(BINARY) $(CMD)

test:
	go test ./internal/... -v -count=1

test-all:
	go test ./... -v -count=1

clean:
	rm -f $(BINARY)

lint:
	go vet ./...

IMAGE ?= nono-nri:latest
KIND_CLUSTER ?= nono-test

docker-build:
	@test -f nono || (echo "ERROR: ./nono binary not found in build context. Download from nono releases and place at repo root." && exit 1)
	docker build -t $(IMAGE) .

docker-load-kind: docker-build
	kind load docker-image $(IMAGE) --name $(KIND_CLUSTER)
