# VARS
BINDIR=./bin
WASMEXEC=./bin/wasm_exec.js

# preparing
checkbindir:
	@if [ ! -d "${BINDIR}" ]; then mkdir "${BINDIR}"; fi

checkwasmexec:
	@if [ ! -d "${WASMEXEC}" ]; then cp $(go env GOROOT)/misc/wasm/wasm_exec.js ./bin; fi

# common build
build: checkbindir
	go build -o ./bin ./cmd/fireplace

run: build
	./bin/fireplace

# wasm build
buildwasm: checkbindir checkwasmexec
	GOOS=js GOARCH=wasm go build -o ./bin ./cmd/web-fireplace

#wasmrun: buildwasm
#	./bin/

# debug
debugbuild:
	go build -tags=ebitenginedebug -o ./bin/fireplace-debug ./cmd/fireplace

debugrun: debugbuild
	./bin/fireplace-debug