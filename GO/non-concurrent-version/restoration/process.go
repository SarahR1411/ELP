package restoration

import (
	"image"
	"image/color"
)

// GetBlendedColorWithEdges calculates a blended color for a pixel at (x, y) by considering the surrounding pixels.
// It uses the pixel's edge information (from the 'edges' matrix) and the weight of the surrounding pixels to compute a new color.
func GetBlendedColorWithEdges(img image.Image, mask [][]float64, edges [][]float64, x, y int) color.Color {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Variables to accumulate the weighted color components and total weight
	var sumR, sumG, sumB, weightSum float64

	// Define the maximum radius for neighboring pixels
	maxRadius := 5

	// Iterate over a square region around the pixel (x, y)
	for dy := -maxRadius; dy <= maxRadius; dy++ {
		for dx := -maxRadius; dx <= maxRadius; dx++ {
			nx, ny := x+dx, y+dy

			// Check if the neighboring pixel is within bounds and not in the mask (mask value < 1.0 indicates valid pixel)
			if nx >= 0 && nx < width && ny >= 0 && ny < height && mask[ny][nx] < 1.0 {
				c := img.At(nx, ny)
				r, g, b, _ := c.RGBA()

				// Calculate the edge weight (based on edge strength)
				edgeWeight := 1.0 - edges[ny][nx]

				// Calculate the distance from the center pixel (x, y)
				distance := float64(dx*dx + dy*dy)

				// Calculate the weight for the neighboring pixel (decays with distance)
				weight := edgeWeight / (1.0 + distance)

				// Accumulate weighted color components and weight
				sumR += float64(r) * weight
				sumG += float64(g) * weight
				sumB += float64(b) * weight
				weightSum += weight
			}
		}
	}

	// If no valid neighbors were found (i.e., weightSum is 0), return the original pixel color
	if weightSum == 0 {
		original := img.At(x, y)
		r, g, b, _ := original.RGBA()
		return color.RGBA{R: uint8(r / 256), G: uint8(g / 256), B: uint8(b / 256), A: 255}
	}

	// Return the weighted average color
	return color.RGBA{
		R: uint8((sumR / weightSum) / 256),
		G: uint8((sumG / weightSum) / 256),
		B: uint8((sumB / weightSum) / 256),
		A: 255,
	}
}


// InpaintWithEdges inpaints an image by filling in the missing areas (marked by the mask) using the blended color from surrounding pixels.
// It takes into account edge information for more natural-looking inpainting.
func InpaintWithEdges(img image.Image, mask [][]float64, edges [][]float64) *image.RGBA {
	bounds := img.Bounds()
	output := image.NewRGBA(bounds)

	// Loop through every pixel in the image bounds
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the feathered mask value for the current pixel
			weight := mask[y][x] // Use feathered mask value
			// If the mask indicates a damaged area (weight > 0), inpaint the pixel
			if weight > 0 {      // If the mask indicates a damaged area
				blendedColor := GetBlendedColorWithEdges(img, mask, edges, x, y)
				output.Set(x, y, blendedColor)
			} else {
				// If the pixel is not in the mask, retain the original pixel
				output.Set(x, y, img.At(x, y)) // Copy original pixel
			}
		}
	}
	// Return the inpainted image
	return output
}