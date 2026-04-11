package main

import (
	"os"
	"time"
)

var (
	b [1760]byte
	z [1760]int8
)

const (
	R1 = 1
	R2 = 2048
	K2 = 5120 * 1024
)

func rot(mul, shift int, x, y *int) {
	t := *x
	*x -= (mul * *y) >> shift
	*y += (mul * t) >> shift
	t = (3145728 - *x**x - *y**y) >> 11
	*x = (*x * t) >> 10
	*y = (*y * t) >> 10
}

func main() {
	w := os.Stdout
	const ramp = ".,-~:;=!*#$@"
	sA, cA := 1024, 0
	sB, cB := 1024, 0
	var out [5 + 1761]byte
	first := true
	for {
		for i := range b {
			b[i] = ' '
			z[i] = 127
		}
		sj, cj := 0, 1024
		for range 90 {
			si, ci := 0, 1024
			for range 324 {
				x0 := R1*cj + R2
				x1 := (ci * x0) >> 10
				x2 := (cA * sj) >> 10
				x3 := (si * x0) >> 10
				x4 := R1*x2 - ((sA * x3) >> 10)
				x5 := (sA * sj) >> 10
				x6 := K2 + R1*1024*x5 + cA*x3
				x7 := (cj * si) >> 10

				x := 40 + 30*(cB*x1-sB*x4)/x6
				y := 12 + 15*(cB*x4+sB*x1)/x6
				n := (((-cA*x7 - cB*(((-sA*x7)>>10)+x2) - ci*((cj*sB)>>10)) >> 10) - x5) >> 7

				o := x + 80*y
				zz := int8((x6 - K2) >> 15)

				if y > 0 && y < 22 && x > 0 && x < 80 && zz < z[o] {
					z[o] = zz
					if n > 0 {
						b[o] = ramp[n]
					} else {
						b[o] = ramp[0]
					}
				}
				rot(5, 8, &ci, &si)
			}
			rot(9, 7, &cj, &sj)
		}
		p := 0
		if !first {
			copy(out[p:], "\x1b[23A")
			p += 5
		}
		for k := range 1761 {
			if k%80 != 0 {
				out[p] = b[k]
			} else {
				out[p] = '\n'
			}
			p++
		}
		_, _ = w.Write(out[:p])
		first = false
		rot(5, 7, &cA, &sA)
		rot(5, 8, &cB, &sB)
		time.Sleep(15 * time.Millisecond)
	}
}
