BINARY_NAME=redis-cli-app
DOCKER_IMAGE=redis-cli-app

GO_BUILD_FLAGS=-o $(BINARY_NAME)

all: build

build:
	go build $(GO_BUILD_FLAGS)

clean:
	rm -f $(BINARY_NAME)

run: build
	@trap 'make clean' EXIT; ./$(BINARY_NAME)

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run -it --rm $(DOCKER_IMAGE)

docker-up: docker-build docker-run

docker-clean:
	docker rmi $(DOCKER_IMAGE)

help:
	@echo "Makefile commands:"
	@echo "  make all          - Build the Go binary (default)"
	@echo "  make build        - Build the Go binary"
	@echo "  make clean        - Remove the built binary"
	@echo "  make run          - Run the Go binary locally"
	@echo "  make docker-build - Build the Docker image"
	@echo "  make docker-run   - Run the Docker container"
	@echo "  make docker-up    - Build and run with Docker"
	@echo "  make docker-clean - Remove the Docker image"
	@echo "  make help         - Display available commands"
