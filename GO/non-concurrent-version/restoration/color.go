package restoration

import (
	"image"
	"image/color"
)


// HistEqual applies histogram equalization to an image for color correction.
// It processes each RGB channel (Red, Green, Blue) separately to enhance image contrast.
func HistEqual(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImg := image.NewRGBA(bounds)

	// Create histograms for each RGB channel (Red, Green, Blue)
	histR, histG, histB := make([]int, 256), make([]int, 256), make([]int, 256)

	// Calculate the histogram for each channel by iterating over all pixels
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			histR[r>>8]++	// Store the frequency of each red component (scaled to [0-255])
			histG[g>>8]++	// Store the frequency of each green component (scaled to [0-255])
			histB[b>>8]++	// Store the frequency of each blue component (scaled to [0-255])
		}
	}

	// Compute the cumulative distribution function (CDF) for each color channel
	cdfR := computeCDF(histR)
	cdfG := computeCDF(histG)
	cdfB := computeCDF(histB)

	// Find min and max CDF values for each channel to normalize the histogram
	minR, maxR := findMinMax(cdfR)
	minG, maxG := findMinMax(cdfG)
	minB, maxB := findMinMax(cdfB)

	// Apply histogram equalization to each pixel in the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the original RGBA values of the current pixel
			r, g, b, a := img.At(x, y).RGBA()

			// Normalize the color values using the CDF
			newR := uint8(((cdfR[r>>8] - minR) * 255) / (maxR - minR))
			newG := uint8(((cdfG[g>>8] - minG) * 255) / (maxG - minG))
			newB := uint8(((cdfB[b>>8] - minB) * 255) / (maxB - minB))

			// Set the new pixel color (with the original alpha value)
			newImg.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)})
		}
	}

	// Return the new image with enhanced color distribution
	return newImg
}

// findMinMax finds the minimum and maximum non-zero values in the CDF array.
func findMinMax(cdf []int) (min, max int) {
	min, max = -1, -1
	for _, value := range cdf {
		if value > 0 && min == -1 {
			min = value	// First non-zero value is the minimum
		}
		if value > 0 {
			max = value	// Last non-zero value is the maximum
		}
	}
	return min, max
}


// computeCDF calculates the cumulative distribution function (CDF) from a given histogram.
// The CDF is used to transform pixel values for histogram equalization.
func computeCDF(hist []int) []int {
	cdf := make([]int, len(hist))
	cdf[0] = hist[0]
	// Compute cumulative sum of the histogram
	for i := 1; i < len(hist); i++ {
		cdf[i] = cdf[i-1] + hist[i]
	}
	return cdf
}

// GetGlobalAverageColor computes the average color (RGB) of an image by averaging all pixels' RGB values.
func GetGlobalAverageColor(img image.Image) color.Color {
	var rSum, gSum, bSum uint64
	var count uint64

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Iterate over every pixel to compute the sum of its RGB components
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			rSum += uint64(r)
			gSum += uint64(g)
			bSum += uint64(b)
			count++
		}
	}

	// Compute average color by dividing the sum of each component by the number of pixels
	// (Note: >> 8 is to convert the value back to [0-255] range)
	return color.RGBA{
		R: uint8((rSum / count) >> 8),
		G: uint8((gSum / count) >> 8),
		B: uint8((bSum / count) >> 8),
		A: 255,
	}
}