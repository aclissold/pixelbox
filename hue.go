package main

import (
	"code.google.com/p/biogo.graphics/palette"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
	size := 1024
	icon := image.NewRGBA(image.Rect(0, 0, size, size))
	hue := &palette.HSVA{0, 1, 1, 1}

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			hue.H = float64(x) / float64(size)
			icon.Set(x, y, hue)
		}
	}

	file, err := os.Create("hue.png")
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(file, icon)
	if err != nil {
		log.Fatal(err)
	}
}
