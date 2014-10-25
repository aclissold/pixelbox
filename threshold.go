package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
)

var in image.Image
var out *image.Gray
var filename string
var bounds image.Rectangle

var wg sync.WaitGroup
var lastUpdate int = 0
var l sync.Mutex

func main() {
	setup()
	draw()
	write()
}

// setup instantiates in, out, filename, and bounds.
func setup() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	usage := "Usage:\n\tgo run blur.go file.png\n"
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening", filename)
		fmt.Println(usage)
		os.Exit(1)
	}
	in, err = png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	bounds = in.Bounds()

	out = image.NewGray(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
}

// draw computes the values for out.
func draw() {
	wg.Add(4)
	go drawquadrant(0, 0) // upper left
	go drawquadrant(0, 1) // upper right
	go drawquadrant(1, 0) // lower left
	go drawquadrant(1, 1) // lower right
	wg.Wait()
	fmt.Println()
}

// drawquadrant computes the values for the (i, j) quadrant of out.
func drawquadrant(j, i int) {
	defer wg.Done()
	lastPercent := -1
	startY, endY := j*bounds.Max.Y/2, bounds.Max.Y/2+j*bounds.Max.Y/2
	startX, endX := i*bounds.Max.X/2, bounds.Max.X/2+i*bounds.Max.X/2
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			bw := threshold(in, x, y, math.MaxUint16/2)
			out.Set(x, y, bw)
		}
		percent := int(25 * (float32(y-startY) / float32(bounds.Max.Y/2)))
		if percent != lastPercent {
			lastPercent = percent
			l.Lock()
			lastUpdate++
			l.Unlock()
			fmt.Printf("\r%d%%", lastUpdate)
		}
	}
}

// write writes out to threshold.png.
func write() {
	file, err := os.Create("threshold.png")
	if err != nil {
		log.Fatal(err)
	}

	if err = png.Encode(file, out); err != nil {
		log.Fatal(err)
	}
}

// threshold returns black if i at (x, y) is <= t; white otherwise.
func threshold(i image.Image, x, y int, t uint32) color.Gray16 {
	gray, _, _, _ := i.At(x, y).RGBA()
	if gray <= t {
		return color.Gray16{0}
	}
	return color.Gray16{math.MaxUint16}
}
