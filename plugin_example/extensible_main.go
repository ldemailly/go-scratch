package main

import (
	"log"
	"plugin"
)

func main() {
	p, err := plugin.Open("plugin/plugin.so")
	if err != nil {
		log.Fatal(err)
	}

	sym, err := p.Lookup("SayHello")
	if err != nil {
		log.Fatal(err)
	}

	sayHello, ok := sym.(func(string))
	if !ok {
		log.Fatal("Invalid function signature")
	}

	sayHello("World")
}
