repro.wasm: wanix.go
	GOOS=wasip1 GOARCH=wasm go build -o repro.wasm -trimpath -ldflags="-w -s" .
	ls -l
