package restoration

import (
	"image"
	"math"
	"runtime"
	"sync"
)

func EdgeDetectionConcurrent(img image.Image) [][]float64 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	edges := make([][]float64, height)
	for i := range edges {
		edges[i] = make([]float64, width)
	}

	// Sobel kernels
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

	// Worker function to process a chunk of pixels
	processChunk := func(xStart, xEnd, yStart, yEnd int, wg *sync.WaitGroup) {
		defer wg.Done()
		for y := yStart; y < yEnd; y++ {
			if y <= 0 || y >= height-1 {
				continue // Skip edge rows
			}
			for x := xStart; x < xEnd; x++ {
				if x <= 0 || x >= width-1 {
					continue // Skip edge columns
				}
				var gx, gy float64
				for ky := -1; ky <= 1; ky++ {
					for kx := -1; kx <= 1; kx++ {
						px := img.At(x+kx, y+ky)
						r, g, b, _ := px.RGBA()
						gray := float64(r+g+b) / (3.0 * 256.0)
						gx += gray * float64(sobelX[ky+1][kx+1])
						gy += gray * float64(sobelY[ky+1][kx+1])
					}
				}
				gradient := math.Sqrt(gx*gx + gy*gy)
				edges[y][x] = gradient

				// Safely update maxGradient
				maxGradientMutex.Lock()
				if gradient > maxGradient {
					maxGradient = gradient
				}
				maxGradientMutex.Unlock()
			}
		}
	}

	// Divide work into chunks
	numWorkers := runtime.NumCPU()
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

	// Normalize and apply threshold
	threshold := 0.2 // You may adjust this value based on testing
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			edges[y][x] /= maxGradient
			if edges[y][x] < threshold {
				edges[y][x] = 0.0
			}
		}
	}

	return edges
}

// FeatherMaskConcurrent applies feathering to a binary mask with concurrency.
func FeatherMaskConcurrent(mask [][]float64, radius int, edgeMask [][]float64) [][]float64 {
	height := len(mask)
	width := len(mask[0])

	// Output mask
	featheredMask := make([][]float64, height)
	for i := range featheredMask {
		featheredMask[i] = make([]float64, width)
	}

	var wg sync.WaitGroup
	numWorkers := 4
	rowsPerWorker := height / numWorkers

	// Feathering worker
	processChunk := func(startRow, endRow int) {
		defer wg.Done()
		for y := startRow; y < endRow; y++ {
			for x := 0; x < width; x++ {
				if mask[y][x] == 1 {
					// Feathering logic
					featheredMask[y][x] = 1.0 // Fully masked
				} else if mask[y][x] == 0 {
					// Unmasked, keep it 0
					featheredMask[y][x] = 0.0
				} else {
					// Feather edges using Gaussian-like falloff
					minDistance := float64(radius)
					for dy := -radius; dy <= radius; dy++ {
						for dx := -radius; dx <= radius; dx++ {
							nx, ny := x+dx, y+dy
							if nx >= 0 && nx < width && ny >= 0 && ny < height && mask[ny][nx] == 1 {
								dist := math.Sqrt(float64(dx*dx + dy*dy))
								if dist < minDistance {
									minDistance = dist
								}
							}
						}
					}
					featheredMask[y][x] = math.Exp(-minDistance / float64(radius))
				}
			}
		}
	}

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		startRow := i * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if i == numWorkers-1 {
			endRow = height
		}
		wg.Add(1)
		go processChunk(startRow, endRow)
	}

	wg.Wait()
	return featheredMask
}
