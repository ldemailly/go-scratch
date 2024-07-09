package main

import (
	"encoding/json"
	"fmt"
)

type A struct {
	Foo string `json:"foo,omitempty"`
}

type B struct {
	Bar string `json:"bar,omitempty"`
}

type Response struct {
	A
	B
}

func main() {
	a := A{Foo: "foo val"}
	b := B{Bar: "bar var"}
	res := Response{}
	res.A = a
	j, _ := json.Marshal(res)
	fmt.Println(string(j))
	res = Response{}
	res.B = b
	j, _ = json.Marshal(res)
	fmt.Println(string(j))
}
