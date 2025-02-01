package restoration

import (
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"sync"
)

// CreateMaskByChunks generates a binary mask of the image using parallel processing.
// It divides the image into chunks and applies a threshold to classify pixels as part of the mask.
// The result is saved as a JPEG image for debugging.
func CreateMaskByChunks(img image.Image, outputPath string, numWorkers int) ([][]float64, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create the mask
	mask := make([][]float64, height)
	for i := range mask {
		mask[i] = make([]float64, width)
	}
	// Define chunk sizes for parallel processing
	chunkWidth := int(math.Ceil(float64(width) / math.Sqrt(float64(numWorkers))))
	chunkHeight := int(math.Ceil(float64(height) / math.Sqrt(float64(numWorkers))))
	overlap := 2	// Overlap to prevent boundary artifacts

	var wg sync.WaitGroup
	processChunk := func(xStart, xEnd, yStart, yEnd int) {
		defer wg.Done()
		for y := yStart; y < yEnd && y < height; y++ {
			for x := xStart; x < xEnd && x < width; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
				sum := uint32(r8) + uint32(g8) + uint32(b8)
				
				// Apply threshold to determine mask value
				if sum > 427 { 
					mask[y][x] = 1.0
				} else {
					mask[y][x] = 0.0
				}
			}
		}
	}

	// Launch workers to process image chunks in parallel
	for yStart := 0; yStart < height; yStart += chunkHeight - overlap {
		yEnd := yStart + chunkHeight
		if yEnd > height {
			yEnd = height
		}
		for xStart := 0; xStart < width; xStart += chunkWidth - overlap {
			xEnd := xStart + chunkWidth
			if xEnd > width {
				xEnd = width
			}
			wg.Add(1)
			go processChunk(xStart, xEnd, yStart, yEnd)
		}
	}
	wg.Wait()

	// Save the mask as a grayscale image for debugging
	maskImg := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if mask[y][x] == 1.0 {
				maskImg.Set(x, y, color.White)
			} else {
				maskImg.Set(x, y, color.Black)
			}
		}
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()
	err = jpeg.Encode(outputFile, maskImg, nil)
	if err != nil {
		return nil, err
	}

	return mask, nil
}

// FeatherMaskConcurrent smooths the edges of a binary mask using an exponential decay function.
// The function runs in parallel, ensuring efficient feathering.
func FeatherMaskConcurrent(mask [][]float64, radius int, edgeMask [][]float64, numWorkers int) [][]float64 {
	height := len(mask)
	width := len(mask[0])

	// Output mask with feathering applied
	featheredMask := make([][]float64, height)
	for i := range featheredMask {
		featheredMask[i] = make([]float64, width)
	}

	var wg sync.WaitGroup
	rowsPerWorker := height / numWorkers

	// Worker function to feather mask in parallel
	processChunk := func(startRow, endRow int) {
		defer wg.Done()
		for y := startRow; y < endRow; y++ {
			for x := 0; x < width; x++ {
				if mask[y][x] == 1 {
					featheredMask[y][x] = 1.0 // Fully masked
				} else {
					for dy := -radius; dy <= radius; dy++ {
						for dx := -radius; dx <= radius; dx++ {
							nx, ny := x+dx, y+dy
							if nx >= 0 && nx < width && ny >= 0 && ny < height && mask[ny][nx] == 1 {
								distance := float64(dx*dx + dy*dy)
								weight := math.Exp(-distance / float64(radius*radius)) * (1.0 - edgeMask[ny][nx])
								featheredMask[y][x] = math.Max(featheredMask[y][x], weight)
							}
						}
					}
				}
			}
		}
	}

	// Launch workers to process feathering
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
