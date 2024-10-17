.PHONY: start
start:
	go build -o ./bin/redis-go-cli 
	./bin/redis-go-cli

.PHONY: test
test:
	go test -v