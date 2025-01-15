package restoration

import (
	"image"
	"image/color"
)

// We'll use Histogram equalization for color correction
func HistEqual(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImg := image.NewRGBA(bounds)

	// Create histograms for each channel
	histR, histG, histB := make([]int, 256), make([]int, 256), make([]int, 256)

	// Calculate the histogram for each channel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			histR[r>>8]++
			histG[g>>8]++
			histB[b>>8]++
		}
	}

	// Calculate the cumulative distribution function (CDF)
	cdfR := computeCDF(histR)
	cdfG := computeCDF(histG)
	cdfB := computeCDF(histB)

	// Find min and max CDF values for normalization
	minR, maxR := findMinMax(cdfR)
	minG, maxG := findMinMax(cdfG)
	minB, maxB := findMinMax(cdfB)

	// Apply histogram equalization to each pixel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			newR := uint8(((cdfR[r>>8] - minR) * 255) / (maxR - minR))
			newG := uint8(((cdfG[g>>8] - minG) * 255) / (maxG - minG))
			newB := uint8(((cdfB[b>>8] - minB) * 255) / (maxB - minB))
			newImg.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)})
		}
	}

	return newImg
}

func findMinMax(cdf []int) (min, max int) {
	min, max = -1, -1
	for _, value := range cdf {
		if value > 0 && min == -1 {
			min = value
		}
		if value > 0 {
			max = value
		}
	}
	return min, max
}


// computeCDF calculates the cumulative distribution function from a histogram
func computeCDF(hist []int) []int {
	cdf := make([]int, len(hist))
	cdf[0] = hist[0]
	for i := 1; i < len(hist); i++ {
		cdf[i] = cdf[i-1] + hist[i]
	}
	return cdf
}

func GetGlobalAverageColor(img image.Image) color.Color {
	var rSum, gSum, bSum uint64
	var count uint64

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

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

	// Compute average
	return color.RGBA{
		R: uint8((rSum / count) >> 8),
		G: uint8((gSum / count) >> 8),
		B: uint8((bSum / count) >> 8),
		A: 255,
	}
}