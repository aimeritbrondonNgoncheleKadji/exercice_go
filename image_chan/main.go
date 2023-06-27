package main

import (
	"image"
	"image/color"
	_ "image/png"
	"image/png"
	"os"
	"time"

	"github.com/disintegration/imaging"
)

func ConvertToGrayScale(input <-chan image.Image, output chan<- image.Image) {
	img := <-input
	bounds := img.Bounds()
	gray := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			gray.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}
	output <- gray
	close(output)
}

func Blur(input <-chan image.Image, output chan<- image.Image) {
	img := <-input
	blurred := imaging.Blur(img, 5)
	output <- blurred
	close(output)
}

func main() {
	src, err := os.Open("src.png")
	if err != nil {
		panic(err)
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		panic(err)
	}

	toGrayscale := make(chan image.Image, 1)
	toGrayscale <- img

	fromGrayscale := make(chan image.Image, 1)
	go ConvertToGrayScale(toGrayscale, fromGrayscale)

	start := time.Now()
	gray := <-fromGrayscale
	elapsed := time.Since(start)
	println("Grayscale conversion took", elapsed.String())

	toBlur := make(chan image.Image, 1)
	toBlur <- gray

	fromBlur := make(chan image.Image, 1)
	go Blur(toBlur, fromBlur)

	start = time.Now()
	blurred := <-fromBlur
	elapsed = time.Since(start)
	println("Blurring took", elapsed.String())

	out, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, blurred)
	if err != nil {
		panic(err)
	}
}
