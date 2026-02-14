package main

import (
	"image"
	"image/color"

	"fortio.org/terminal/ansipixels"
)

func heartEquation(x, y float64) float64 {
	a := x*x + y*y - 1
	return a*a*a - x*x*y*y*y
}

func main() {
	ap := ansipixels.NewAnsiPixels(0)
	if err := ap.Open(); err != nil {
		panic(err)
	}
	defer ap.Restore()
	ap.SyncBackgroundColor()
	size := ap.W
	if 2*ap.H < ap.W {
		size = 2 * ap.H
	}
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for py := 0; py < size; py++ {
		for px := 0; px < size; px++ {
			// Convert pixel coords to mathematical coords
			x := -1.5 + 3.0*float64(px)/float64(size-1)
			y := 1.5 - 3.0*float64(py)/float64(size-1) // Flip y for screen coords
			z := heartEquation(x, y)
			if z <= 0 {
				img.Set(px, py, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			}
		}
	}
	ap.ClearScreen()
	ap.ShowScaledImage(img)
	msg := "Happy Valentine [(x^2+y^2-1)^3 = x^2 y^3]"
	ap.WriteAt(size/2-len(msg)/2, 0, msg)
	ap.MoveCursor(0, ap.H-1)
}
