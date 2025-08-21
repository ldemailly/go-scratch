GOOS?=wasip1
GO=go
GOARCH?=wasm

all: clean repro.wasm

repro.wasm: wanix.go
	# GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o repro.wasm -trimpath -ldflags="-w -s" .
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o repro.wasm .
	ls -l

clean:
	rm -f repro.wasm

.PHONY: all clean
