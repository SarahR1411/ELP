package restoration

import (
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"sync"
	"runtime"
)






// These functions need to be improved. They're still kinda slow and they create weird lines in the image 








func CreateMaskByChunks(img image.Image, outputPath string, numWorkers int) ([][]float64, error) {
	// Get the dimensions of the image
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// create a new mask (2D array of floats) where 1.0 represents a 'damaged' pixel and 0.0 a normal one
	mask := make([][]float64, height)
	for i:= range mask{
		mask[i] = make([]float64, height)
	}
	

	rowsPerWorker := height/numWorkers
	var wg sync.WaitGroup

	processChunk := func(startRow, endRow int) {
		defer wg.Done()
		for y := startRow; y < endRow; y++ {
			for x := 0; x < width; x++ {
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

func FeatherMaskConcurrent(mask [][]float64, radius int, edgeMask [][]float64) [][]float64 {
    height := len(mask)
    width := len(mask[0])
    feathered := make([][]float64, height)
    for y := range feathered {
        feathered[y] = make([]float64, width)
    }

    threshold := 0.01 // value needs testing

    // Worker function to process rows with overlap
    processRow := func(yStart, yEnd int, wg *sync.WaitGroup) {
        defer wg.Done()
        for y := yStart; y < yEnd; y++ {
            for x := 0; x < width; x++ {
                if mask[y][x] > 0 {
                    feathered[y][x] = 1.0
                    for dy := -radius; dy <= radius; dy++ {
                        for dx := -radius; dx <= radius; dx++ {
                            nx, ny := x+dx, y+dy
                            if nx >= 0 && nx < width && ny >= 0 && ny < height {
                                distance := float64(dx*dx + dy*dy)
                                weight := math.Exp(-distance / float64(radius*radius))
                                if ny < len(edgeMask) && nx < len(edgeMask[ny]) { // Ensure edgeMask bounds
                                    weight *= 1.0 - edgeMask[ny][nx]
                                }
                                // Apply threshold to prevent negligible weights
                                if weight > threshold {
                                    feathered[ny][nx] = math.Max(feathered[ny][nx], weight)
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    // Divide work among workers with overlap
    numWorkers := runtime.NumCPU()
    wg := &sync.WaitGroup{}
    rowsPerWorker := height / numWorkers
    for i := 0; i < numWorkers; i++ {
        yStart := i * rowsPerWorker
        yEnd := (i + 1) * rowsPerWorker
        if i == numWorkers-1 { // Handle remaining rows for the last worker
            yEnd = height
        }

        // Extend the range to include overlap
        yStartWithOverlap := max(0, yStart-radius)
        yEndWithOverlap := min(height, yEnd+radius)

        wg.Add(1)
        go processRow(yStartWithOverlap, yEndWithOverlap, wg)
    }
    wg.Wait()

    return feathered
}

// Helper functions to clamp values
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
 

func SaveFeatheredMask(feathered [][]float64, outputPath string) error {
    height := len(feathered)
    width := len(feathered[0])
    featheredImg := image.NewGray(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            intensity := uint8(feathered[y][x] * 255) // Normalize to [0, 255]
            featheredImg.SetGray(x, y, color.Gray{Y: intensity})
        }
    }
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()
    return jpeg.Encode(file, featheredImg, nil)
}



