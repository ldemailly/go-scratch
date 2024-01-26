
run-demo: demo
	./demo

demo: *.go plugin/plugin.so
	go build -o demo .
	go version -m ./demo


plugin/plugin.so:
	make -C plugin

clean:
	rm demo plugin/plugin.so
