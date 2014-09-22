package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

var size int
var img *image.RGBA
var r *rand.Rand

func main() {
	setup()
	draw()
	write()
}

func setup() {
	size = 1024
	img = image.NewRGBA(image.Rect(0, 0, size, size))

	t := time.Now().UnixNano()
	s := rand.NewSource(t)
	r = rand.New(s)
}

func draw() {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			gray := color.Gray16{uint16(r.Intn(65536))}
			img.Set(x, y, gray)
		}
	}
}

func write() {
	file, err := os.Create("noise.png")
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}
}
