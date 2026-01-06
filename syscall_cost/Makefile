

all: bench


ccat: ccat.c
	gcc -Wall -Wextra -pedantic -Werror -O3 -o ccat ccat.c

gocat: gocat.go
	go build -o gocat gocat.go


/tmp/bigfile:
	dd if=/dev/zero of=/tmp/bigfile bs=1M count=100

bench: /tmp/bigfile gocat ccat
	hyperfine --warmup=5 "./ccat < /tmp/bigfile > /dev/null" "./gocat < /tmp/bigfile > /dev/null"

.PHONY: bench
