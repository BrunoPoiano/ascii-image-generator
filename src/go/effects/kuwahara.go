package effects

import (
	"image"
	"image/color"
	"math"
)

func clamp(value, minVal, maxVal int) int {
	return int(math.Max(float64(minVal), math.Min(float64(value), float64(maxVal))))
}

func rgbToHSV(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	rf := float64(r) / 65535.0
	gf := float64(g) / 65535.0
	bf := float64(b) / 65535.0

	max := math.Max(math.Max(rf, gf), bf)
	return max
}

func stdDev(values []float64) float64 {
	n := float64(len(values))
	if n == 0 {
		return 0
	}

	// Mean calculation
	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / n

	// Standard deviation calculation
	var variance float64
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	return math.Sqrt(variance / n)
}

func averageRGB(img image.Image, x1, y1, x2, y2 int) color.Color {
	var sumR, sumG, sumB float64
	count := 0

	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			sumR += float64(r)
			sumG += float64(g)
			sumB += float64(b)
			count++
		}
	}

	if count == 0 {
		return color.Black
	}

	return color.RGBA{
		R: uint8(sumR / float64(count) / 256),
		G: uint8(sumG / float64(count) / 256),
		B: uint8(sumB / float64(count) / 256),
		A: 255,
	}
}

func KuwaharaFilter(img image.Image, size int) image.Image {

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	result := image.NewRGBA(bounds)

	halfWin := size / 2
	quadSize := int(math.Ceil(float64(size) / 2.0))

	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			//quadrants
			tlX, tlY := x-halfWin, y-halfWin
			q1X, q1Y := tlX+quadSize, tlY+quadSize
			q2X, q2Y := tlX+size, tlY+quadSize
			q3X, q3Y := tlX+quadSize, tlY+size
			q4X, q4Y := tlX+size, tlY+size

			//clamp values
			q1X, q1Y = clamp(q1X, 0, width), clamp(q1Y, 0, height)
			q2X, q2Y = clamp(q2X, 0, width), clamp(q2Y, 0, height)
			q3X, q3Y = clamp(q3X, 0, width), clamp(q3Y, 0, height)
			q4X, q4Y = clamp(q4X, 0, width), clamp(q4Y, 0, height)

			//extracting brightness
			quadrands := []struct {
				x1, y1, x2, y2 int
			}{
				{tlX, tlY, q1X, q1Y},
				{q1X, tlY, q2X, q2Y},
				{tlX, q1Y, q3X, q3Y},
				{q1X, q1Y, q4X, q4Y},
			}

			//calc standard deviation
			minStdDev := math.MaxFloat64
			bestQuad := quadrands[0]

			for _, quad := range quadrands {
				var values []float64
				for yy := quad.y1; yy < quad.y2; yy++ {
					for xx := quad.x1; xx < quad.x2; xx++ {
						values = append(values, rgbToHSV(img.At(xx, yy)))
					}
				}

				stdDevVal := stdDev(values)
				if stdDevVal < minStdDev {
					minStdDev = stdDevVal
					bestQuad = quad
				}
			}

			// Assign the average color of the best quadrant
			result.Set(x, y, averageRGB(img, bestQuad.x1, bestQuad.y1, bestQuad.x2, bestQuad.y2))
		}
	}

	return result
}
