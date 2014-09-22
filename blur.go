package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
)

var size int
var in *image.RGBA
var r *rand.Rand

func main() {
	f, err := os.Open("noise.png")
	if err != nil {
		log.Fatal("you must run noise.go first")
	}
	in, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	b := in.Bounds()

	out := image.NewGray(image.Rect(0, 0, b.Max.X, b.Max.Y))
	usage := "Usage:\n\tgo run blur.go r\n\n" +
		"Where r is on the interval [0, 1024).\n"
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	r, err := strconv.Atoi(os.Args[1])
	if err != nil || r < 0 || r > 1023 {
		fmt.Println(usage)
		os.Exit(1)
	}
	lastUpdate := ""
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			gray := mean(in, x, y, r)
			out.Set(x, y, gray)
		}
		percent := 100 * (float32(y) / float32(b.Max.Y))
		update := fmt.Sprintf("\r%.0f%%", percent)
		if update != lastUpdate {
			fmt.Print(update)
			lastUpdate = update
		}
	}
	fmt.Println()

	perm := os.ModeDir | 0755
	if err := os.Mkdir("blur", perm); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	filename := fmt.Sprintf("blur/blur%d.png", r)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(file, out)
	if err != nil {
		log.Fatal(err)
	}
}

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
