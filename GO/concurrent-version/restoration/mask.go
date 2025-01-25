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

func CreateMaskByChunks(img image.Image, outputPath string) ([][]float64, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create the mask
	mask := make([][]float64, height)
	for i := range mask {
		mask[i] = make([]float64, width)
	}

	numWorkers := runtime.NumCPU()
	chunkWidth := int(math.Ceil(float64(width) / math.Sqrt(float64(numWorkers))))
	chunkHeight := int(math.Ceil(float64(height) / math.Sqrt(float64(numWorkers))))
	overlap := 2

	var wg sync.WaitGroup
	processChunk := func(xStart, xEnd, yStart, yEnd int) {
		defer wg.Done()
		for y := yStart; y < yEnd && y < height; y++ {
			for x := xStart; x < xEnd && x < width; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				luminance := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
				if luminance > 180 { // Debug threshold
					mask[y][x] = 1.0
				} else {
					mask[y][x] = 0.0
				}
			}
		}
	}

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

	// Save mask for debugging
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
