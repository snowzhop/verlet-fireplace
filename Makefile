build:
	go build -o ./bin ./cmd/fireplace

run: build
	./bin/fireplace