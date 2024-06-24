package idiv

func Div1(a, b int32) (int32, int32) {
	return a / b, a % b
}

func Div2(a, b int32) (int32, int32) {
	d := a / b
	return d, a - d*b
}

// 'bad' version

func abs1(x int32) int32 {
	if x < 0 {
		x *= -1
	}
	return x
}

func abs2(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func Div3(a, b int32) (int32, int32) {
	d := a / b
	return d, abs1(a) - abs2(d*b)
}
