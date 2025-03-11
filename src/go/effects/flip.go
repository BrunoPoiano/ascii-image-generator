package effects

import (
	"image"
)

func FlipHorizontal(img image.Image) image.Image {

	bounds := img.Bounds()
	newImage := image.NewNRGBA64(bounds)

	for y := 0; y < bounds.Max.Y; y++ {
		newX := bounds.Max.X - 1
		for x := 0; x < bounds.Max.X; x++ {
			pixel := img.At(newX, y)
			newImage.Set(x, y, pixel)
			newX--
		}
	}

	return newImage
}

func FlipVertical(img image.Image) image.Image {

	bounds := img.Bounds()
	newImage := image.NewNRGBA64(bounds)

	for x := 0; x < bounds.Max.X; x++ {
		newY := bounds.Max.Y - 1
		for y := 0; y < bounds.Max.Y; y++ {
			pixel := img.At(x, newY)
			newImage.Set(x, y, pixel)
			newY--
		}
	}

	return newImage
}
