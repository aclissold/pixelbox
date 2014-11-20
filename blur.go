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

var numcpu int

var wg sync.WaitGroup
var sectionUpdates float32 = 0
var lastPercent int = -1
var maxSectionUpdates float32
var l sync.Mutex

func main() {
	setup()
	draw()
	write()
}

// setup instantiates in, out, radius, and bounds.
func setup() {
	numcpu = runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)

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
	maxSectionUpdates = float32(bounds.Max.Y)
}

// draw computes the values for out.
func draw() {
	wg.Add(numcpu)
	for i := 0; i < numcpu; i++ {
		go drawsection(i)
	}
	wg.Wait()
	fmt.Println()
}

// drawsection computes the values for the ith partition of out.
func drawsection(i int) {
	defer wg.Done()

	width := bounds.Max.X / numcpu
	startX := i * width
	endX := startX + width
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < endX; x++ {
			gray := mean(in, x, y, radius)
			out.Set(x, y, gray)
		}

		l.Lock()
		sectionUpdates++
		l.Unlock()
		currentPercent := int((100 * (sectionUpdates / maxSectionUpdates)) / float32(numcpu))
		if currentPercent != lastPercent {
			lastPercent = currentPercent
			fmt.Printf("\r%d%%", lastPercent)
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
