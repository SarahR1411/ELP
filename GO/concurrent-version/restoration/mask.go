package restoration

import (
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"runtime"
	"sync"
)

// These functions need to be improved. They're still kinda slow and they create weird lines in the image

func CreateMaskByChunks(img image.Image, outputPath string) ([][]float64, error) {
	// Get the dimensions of the image
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create a new mask (2D array of floats)
	mask := make([][]float64, height)
	for i := range mask {
		mask[i] = make([]float64, width)
	}

	// Get the number of CPU cores
	numWorkers := runtime.NumCPU()

	// Dynamically calculate chunk size
	chunkWidth := int(math.Ceil(float64(width) / math.Sqrt(float64(numWorkers))))
	chunkHeight := int(math.Ceil(float64(height) / math.Sqrt(float64(numWorkers))))

	var wg sync.WaitGroup

	// Worker function to process a chunk of the image
	processChunk := func(xStart, xEnd, yStart, yEnd int) {
		defer wg.Done()
		for y := yStart; y < yEnd && y < height; y++ {
			for x := xStart; x < xEnd && x < width; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				sum := uint32(r>>8) + uint32(g>>8) + uint32(b>>8)
				if sum > 427 { // Adjust threshold as needed
					mask[y][x] = 1.0
				} else {
					mask[y][x] = 0.0
				}
			}
		}
	}

	// Launch goroutines for each chunk
	for yStart := 0; yStart < height; yStart += chunkHeight {
		for xStart := 0; xStart < width; xStart += chunkWidth {
			xEnd := xStart + chunkWidth
			yEnd := yStart + chunkHeight
			wg.Add(1)
			go processChunk(xStart, xEnd, yStart, yEnd)
		}
	}

	wg.Wait()

	// Save the mask image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

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
	err = jpeg.Encode(outputFile, maskImg, nil)
	if err != nil {
		return nil, err
	}

	return mask, nil
}
