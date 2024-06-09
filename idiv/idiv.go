package idiv

func Div1(a, b int) (int, int) {
	return a / b, a % b
}

func Div2(a, b int) (int, int) {
	d := a / b
	return d, a - d*b
}
