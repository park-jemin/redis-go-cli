.PHONY: start
start:
	go build -o ./bin/redis-go-cli 
	./bin/redis-go-cli