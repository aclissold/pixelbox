package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
)

var in image.Image
var out *image.Gray
var radius int
var bounds image.Rectangle

func main() {
	setup()
	draw()
	write()
}

// setup instantiates in, out, radius, and bounds.
func setup() {
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
	lastUpdate := ""
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			gray := mean(in, x, y, radius)
			out.Set(x, y, gray)
		}
		percent := 100 * (float32(y) / float32(bounds.Max.Y))
		update := fmt.Sprintf("\r%.0f%%", percent)
		if update != lastUpdate {
			fmt.Print(update)
			lastUpdate = update
		}
	}
	fmt.Println()
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
