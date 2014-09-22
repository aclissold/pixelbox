package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
)

var in image.Image
var out *image.Gray
var radius int
var bounds image.Rectangle

var wg sync.WaitGroup
var lastUpdate int = 0
var l sync.Mutex

func main() {
	setup()
	draw()
	write()
}

// setup instantiates in, out, radius, and bounds.
func setup() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	f, err := os.Open("noise.png")
	if err != nil {
		log.Fatal("you must run noise.go first")
	}
	in, err = png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	bounds = in.Bounds()

	out = image.NewGray(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	usage := "Usage:\n\tgo run blur.go radius\n\n" +
		"Where radius is on the interval [0, 1024).\n"
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	radius, err = strconv.Atoi(os.Args[1])
	if err != nil || radius < 0 || radius > 1023 {
		fmt.Println(usage)
		os.Exit(1)
	}
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
			gray := mean(in, x, y, radius)
			out.Set(x, y, gray)
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

// write writes out to blur/blur<radius>.png.
func write() {
	perm := os.ModeDir | 0755
	if err := os.Mkdir("blur", perm); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	filename := fmt.Sprintf("blur/blur%d.png", radius)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	if err = png.Encode(file, out); err != nil {
		log.Fatal(err)
	}
}

// mean returns the mean value of the square centered at
// (x, y) in i with radius r.
func mean(i image.Image, x, y, r int) color.Gray16 {
	min := image.Point{X: x - r, Y: y - r}
	max := image.Point{X: x + r, Y: y + r}
	rect := image.Rectangle{min, max}

	var avg uint32 = 0
	var count uint32 = 0
	for y := rect.Min.Y; y <= rect.Max.Y; y++ {
		for x := rect.Min.X; x <= rect.Max.X; x++ {
			c := i.At(x, y)
			r, _, _, _ := c.RGBA()
			avg = (avg*count + r) / (count + 1)
			count++
		}
	}

	return color.Gray16{uint16(avg)}
}
