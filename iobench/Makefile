
all: test bench

test: iobench
	@echo "--- 3 lines ---"
	./iobench -mode scanner -filename 3lines.txt
	@echo ""
	./iobench -mode optio -filename 3lines.txt
	@echo ""
	@echo "--- 2 lines not eol at the end of last ---"
	./iobench -mode scanner -filename 2lines_no_last_eol.txt
	@echo ""
	./iobench -mode optio -filename 2lines_no_last_eol.txt


iobench: Makefile main.go */*.go
	go build .

bench: genfile/random_lines.txt iobench
	# warm up / read file in:
	./iobench -mode scanner
	time ./iobench -mode scanner
	./iobench -mode optio
	time ./iobench -mode optio
	./iobench -mode blockio
	time ./iobench -mode blockio

genfile/random_lines.txt: genfile/genfile.go
	cd genfile && go run .


clean:
	rm -f iobench
	rm -f genfile/random_lines.txt
