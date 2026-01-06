

all: bench


ccat: ccat.c
	gcc -Wall -Wextra -pedantic -Werror -O3 -o ccat ccat.c

gocat24: gocat.go
	GOTOOLCHAIN=go1.24.11 go build -o gocat24 gocat.go

gocat25: gocat.go
	GOTOOLCHAIN=go1.25.5 go build -o gocat25 gocat.go

gocat26: gocat.go
	GOTOOLCHAIN=go1.26rc1 go build -o gocat26 gocat.go

/tmp/bigfile:
	dd if=/dev/zero of=/tmp/bigfile bs=1M count=100

bench: /tmp/bigfile gocat24 gocat25 gocat26 ccat
	hyperfine --warmup=5 "./ccat < /tmp/bigfile > /dev/null" "./gocat24 < /tmp/bigfile > /dev/null" "./gocat25 < /tmp/bigfile > /dev/null" "./gocat26 < /tmp/bigfile > /dev/null"

.PHONY: bench
