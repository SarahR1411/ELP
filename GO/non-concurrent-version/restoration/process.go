package restoration

import (
	"image"
	"image/color"
)

// Get the median color of surrounding pixels for in painting
func GetBlendedColorWithEdges(img image.Image, mask [][]float64, edges [][]float64, x, y int) color.Color {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var sumR, sumG, sumB, weightSum float64
	maxRadius := 5

	for dy := -maxRadius; dy <= maxRadius; dy++ {
		for dx := -maxRadius; dx <= maxRadius; dx++ {
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < width && ny >= 0 && ny < height && mask[ny][nx] < 1.0 {
				c := img.At(nx, ny)
				r, g, b, _ := c.RGBA()
				edgeWeight := 1.0 - edges[ny][nx]
				distance := float64(dx*dx + dy*dy)
				weight := edgeWeight / (1.0 + distance)
				sumR += float64(r) * weight
				sumG += float64(g) * weight
				sumB += float64(b) * weight
				weightSum += weight
			}
		}
	}

	if weightSum == 0 {
		original := img.At(x, y)
		r, g, b, _ := original.RGBA()
		return color.RGBA{R: uint8(r / 256), G: uint8(g / 256), B: uint8(b / 256), A: 255}
	}

	return color.RGBA{
		R: uint8((sumR / weightSum) / 256),
		G: uint8((sumG / weightSum) / 256),
		B: uint8((sumB / weightSum) / 256),
		A: 255,
	}
}


// Inpaint colors the white lines and stains in an image with the median color of surrounding pixels.
func InpaintWithEdges(img image.Image, mask [][]float64, edges [][]float64) *image.RGBA {
	bounds := img.Bounds()
	output := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			weight := mask[y][x] // Use feathered mask value
			if weight > 0 {      // If the mask indicates a damaged area
				blendedColor := GetBlendedColorWithEdges(img, mask, edges, x, y)
				output.Set(x, y, blendedColor)
			} else {
				output.Set(x, y, img.At(x, y)) // Copy original pixel
			}
		}
	}

	return output
}