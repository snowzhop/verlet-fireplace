build:
	go build -o ./bin ./cmd/fireplace

run: build
	./bin/fireplace

debugbuild:
	go build -tags=ebitenginedebug -o ./bin/fireplace-debug ./cmd/fireplace

debugrun: debugbuild
	./bin/fireplace-debug