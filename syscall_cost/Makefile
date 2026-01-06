

all: bench1 bench2

WARMUP:=5

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

bench1: /tmp/bigfile gocat24 gocat25 gocat26 ccat
	hyperfine --warmup=$(WARMUP) \
		"./ccat < /tmp/bigfile > /dev/null" \
		"./gocat24 -syscall < /tmp/bigfile > /dev/null" \
		"./gocat25 -syscall < /tmp/bigfile > /dev/null" \
		"./gocat26 -syscall < /tmp/bigfile > /dev/null"

bench2: /tmp/bigfile gocat26
	hyperfine --warmup=$(WARMUP) \
		"./gocat26 < /tmp/bigfile > /dev/null" \
		"./gocat26 -syscall < /tmp/bigfile > /dev/null"

clean:
	$(RM) /tmp/bigfile ccat gocat24 gocat25 gocat26

.PHONY: bench1 bench2 clean all
