# VARS
BINDIR=./bin
WASMEXEC=./bin/wasm_exec.js
MAINHTML=./bin/index.html
GOROOT=$(shell go env GOROOT)

.PHONY: add-main-html.sh

# preparing
checkbindir:
	@if [ ! -d "${BINDIR}" ]; then mkdir "${BINDIR}"; fi

checkwasmexec:
	@if [ ! -d "${WASMEXEC}" ]; then cp ${GOROOT}/misc/wasm/wasm_exec.js ./bin; fi

checkmainhtml:
	chmod +x ./add-main-html.sh
	@if [ ! -d "${MAINHTML}" ]; then ./add-main-html.sh; fi

# common build
build: checkbindir
	go build -o ./bin ./cmd/fireplace

run: build
	./bin/fireplace

# wasm build
buildwasm: checkbindir checkwasmexec checkmainhtml
	GOOS=js GOARCH=wasm go build -o ./bin/fireplace.wasm ./cmd/fireplace

#wasmrun: buildwasm
#	./bin/

# debug
debugbuild:
	go build -tags=ebitenginedebug -o ./bin/fireplace-debug ./cmd/fireplace

debugrun: debugbuild
	./bin/fireplace-debug