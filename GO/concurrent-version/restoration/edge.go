package restoration

import (
	"image"
	"math"
	"sync"
)

// EdgeDetectionConcurrent performs Sobel edge detection on an image using concurrent processing.
// It splits the image into chunks processed in parallel, computing gradient magnitudes for each pixel.
func EdgeDetectionConcurrent(img image.Image, numWorkers int) [][]float64 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	edges := make([][]float64, height)
	for i := range edges {
		edges[i] = make([]float64, width)
	}

	// Sobel kernels for gradient computation
	sobelX := [][]int{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	sobelY := [][]int{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	var maxGradient float64
	maxGradientMutex := &sync.Mutex{} // Protects access to maxGradient

	// Worker function to process a chunk of the image
	processChunk := func(xStart, xEnd, yStart, yEnd int, wg *sync.WaitGroup) {
		defer wg.Done()
		for y := yStart; y < yEnd; y++ {
			if y <= 0 || y >= height-1 {
				continue // Skip edge rows to avoid out-of-bounds access
			}
			for x := xStart; x < xEnd; x++ {
				if x <= 0 || x >= width-1 {
					continue // Skip edge columns to avoid out-of-bounds access
				}
				var gx, gy float64
				// Compute gradients using Sobel operator
				for ky := -1; ky <= 1; ky++ {
					for kx := -1; kx <= 1; kx++ {
						px := img.At(x+kx, y+ky)
						r, g, b, _ := px.RGBA()
						gray := float64(r+g+b) / (3.0 * 256.0)	// Convert to grayscale
						gx += gray * float64(sobelX[ky+1][kx+1])
						gy += gray * float64(sobelY[ky+1][kx+1])
					}
				}
				// Compute gradient magnitude
				gradient := math.Sqrt(gx*gx + gy*gy)
				edges[y][x] = gradient

				// Safely update maxGradient to track the highest gradient value
				maxGradientMutex.Lock()
				if gradient > maxGradient {
					maxGradient = gradient
				}
				maxGradientMutex.Unlock()
			}
		}
	}

	// Divide work into chunks and distribute to workers
	chunkSize := int(math.Sqrt(float64((width * height) / numWorkers)))
	wg := &sync.WaitGroup{}
	for yStart := 0; yStart < height; yStart += chunkSize {
		for xStart := 0; xStart < width; xStart += chunkSize {
			xEnd := xStart + chunkSize
			yEnd := yStart + chunkSize
			if xEnd > width {
				xEnd = width
			}
			if yEnd > height {
				yEnd = height
			}
			wg.Add(1)
			go processChunk(xStart, xEnd, yStart, yEnd, wg)
		}
	}
	wg.Wait()

	// Normalize the gradient values and apply a threshold for edge detection
	threshold := 0.2
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			edges[y][x] /= maxGradient	// Normalize gradient values
			if edges[y][x] < threshold {
				edges[y][x] = 0.0	// Normalize gradient values
			}
		}
	}

	return edges
}


