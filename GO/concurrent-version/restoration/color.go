package restoration

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// HistEqualConcurrent applies histogram equalization to an image using concurrent processing.
// This enhances the contrast of the image by redistributing pixel intensity values.
func HistEqualConcurrent(img image.Image, numWorkers int) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImg := image.NewRGBA(bounds)

	// Create histograms for each channel
	histR, histG, histB := make([]int, 256), make([]int, 256), make([]int, 256)

	// Mutex to synchronize access to shared histogram data
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Determine chunk size for dividing work among workers
	chunkHeight := int(math.Ceil(float64(height) / float64(numWorkers)))

	// Function to process a chunk of the image and calculate local histograms
	processChunk := func(startY, endY int) {
		defer wg.Done()
		localHistR, localHistG, localHistB := make([]int, 256), make([]int, 256), make([]int, 256)

		for y := startY; y < endY; y++ {
			for x := 0; x < width; x++ {
				c := img.At(x, y)
				r, g, b, _ := c.RGBA()
				localHistR[r>>8]++
				localHistG[g>>8]++
				localHistB[b>>8]++
			}
		}

		// Merge local histograms into the global histograms
		mu.Lock()
		for i := 0; i < 256; i++ {
			histR[i] += localHistR[i]
			histG[i] += localHistG[i]
			histB[i] += localHistB[i]
		}
		mu.Unlock()
	}

	// Launch workers to compute histograms in parallel
	for startY := 0; startY < height; startY += chunkHeight {
		endY := startY + chunkHeight
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go processChunk(startY, endY)
	}

	wg.Wait()

	// Compute cumulative distribution function (CDF) for each color channel
	cdfR := computeCDF(histR)
	cdfG := computeCDF(histG)
	cdfB := computeCDF(histB)

	// Find minimum and maximum CDF values for normalization
	minR, maxR := findMinMax(cdfR)
	minG, maxG := findMinMax(cdfG)
	minB, maxB := findMinMax(cdfB)

	// Function to apply histogram equalization to a chunk of the image
	wg = sync.WaitGroup{}
	processEqualizationChunk := func(startY, endY int) {
		defer wg.Done()
		for y := startY; y < endY; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				newR := uint8(((cdfR[r>>8] - minR) * 255) / (maxR - minR))
				newG := uint8(((cdfG[g>>8] - minG) * 255) / (maxG - minG))
				newB := uint8(((cdfB[b>>8] - minB) * 255) / (maxB - minB))
				newImg.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)})
			}
		}
	}

	// Launch workers to apply histogram equalization in parallel
	for startY := 0; startY < height; startY += chunkHeight {
		endY := startY + chunkHeight
		if endY > height {
			endY = height
		}
		wg.Add(1)
		go processEqualizationChunk(startY, endY)
	}

	wg.Wait()
	return newImg
}

// findMinMax finds the minimum and maximum non-zero values in a CDF for normalization
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

// GetGlobalAverageColor calculates the average color of an image by summing all pixel values
func GetGlobalAverageColor(img image.Image) color.Color {
	var rSum, gSum, bSum uint64
	var count uint64

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Iterate over each pixel to sum up RGB values
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

	// Compute the average color
	return color.RGBA{
		R: uint8((rSum / count) >> 8),
		G: uint8((gSum / count) >> 8),
		B: uint8((bSum / count) >> 8),
		A: 255,
	}
}
