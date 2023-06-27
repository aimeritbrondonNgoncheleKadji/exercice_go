package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"sync"

	"github.com/disintegration/imaging"
)

// ConvertToGrayScale converts an image to grayscale.
func ConvertToGrayScale(img image.Image, wg *sync.WaitGroup) image.Image {
	defer wg.Done()

	bounds := img.Bounds()
	gray := image.NewRGBA(bounds)

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			gray.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	return gray
}

// Blur blurs an image.
func Blur(img image.Image, wg *sync.WaitGroup) image.Image {
	defer wg.Done()

	return imaging.Blur(img, 5)
}

func main() {
	// Open the image.
	src, err := os.Open("src.png")
	if err != nil {
		panic(err)
	}
	defer src.Close()

	// Decode the image.
	img, _, err := image.Decode(src)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	// Convert to grayscale.
	wg.Add(1)
	go func() {
		img = ConvertToGrayScale(img, &wg)
	}()

	// Wait for the grayscale conversion to finish.
	wg.Wait()

	// Apply blur.
	wg.Add(1)
	go func() {
		img = Blur(img, &wg)
	}()

	// Wait for the blur to finish.
	wg.Wait()

	// Save the result.
	out, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		panic(err)
	}
}
