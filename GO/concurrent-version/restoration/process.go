package restoration

import (
	"image"
	"image/color"
	"math"
	"sync"
	"runtime"
)

// Get the median color of surrounding pixels for inpainting
func GetBlendedColorWithEdges(img image.Image, mask [][]float64, edges [][]float64, x, y int) color.Color {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var sumR, sumG, sumB, weightSum float64
	maxRadius := 5
	adjustedRadius := maxRadius
	if x < maxRadius {
		adjustedRadius = x
	} else if x >= width-maxRadius {
		adjustedRadius = width - x - 1
	}
	if y < maxRadius {
		adjustedRadius = min(adjustedRadius, y)
	} else if y >= height-maxRadius {
		adjustedRadius = min(adjustedRadius, height - y - 1)
	}

	for dy := -adjustedRadius; dy <= adjustedRadius; dy++ {
		for dx := -adjustedRadius; dx <= adjustedRadius; dx++ {
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < width && ny >= 0 && ny < height && mask[ny][nx] < 1.0 {
				c := img.At(nx, ny)
				r, g, b, _ := c.RGBA()
				edgeWeight := 1.0 - edges[ny][nx]
				distance := float64(dx*dx + dy*dy)
				weight := edgeWeight / (math.Sqrt(distance) + 1e-6)
				sumR += float64(r) * weight
				sumG += float64(g) * weight
				sumB += float64(b) * weight
				weightSum += weight
			}
		}
	}
	

	if weightSum == 0 {
		return img.At(x, y) // Use the original image color directly 
	}
	


	return color.RGBA{
		R: uint8((sumR / weightSum) / 256),
		G: uint8((sumG / weightSum) / 256),
		B: uint8((sumB / weightSum) / 256),
		A: 255,
	}
}

// InpaintByChunks processes the image by dividing it into chunks of pixels and applying inpainting.
func InpaintByChunks(img image.Image, mask [][]float64, edges [][]float64) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	output := image.NewRGBA(bounds)

	numWorkers := runtime.NumCPU() // Use available CPU cores
	chunkWidth := int(math.Ceil(float64(width)/math.Sqrt(float64(numWorkers)))) + 1
	chunkHeight := int(math.Ceil(float64(height)/math.Sqrt(float64(numWorkers)))) + 1


	var wg sync.WaitGroup

	// Worker function for processing a chunk of the image
	overlap := 5 // Number of pixels to overlap
	var mu sync.Mutex
	processChunk := func(xStart, xEnd, yStart, yEnd int) {
    defer wg.Done()
    for y := max(0, yStart-overlap); y < min(yEnd+overlap, height); y++ {
        for x := max(0, xStart-overlap); x < min(xEnd+overlap, width); x++ {
            mu.Lock()
            if mask[y][x] > 0 {
                blendedColor := GetBlendedColorWithEdges(img, mask, edges, x, y)
				output.Set(x, y, blendedColor)
            } else {
                output.Set(x, y, img.At(x, y))
            }
            mu.Unlock()
        }
    }
}

	for yStart := 0; yStart < height; yStart += chunkHeight - overlap {
		for xStart := 0; xStart < width; xStart += chunkWidth - overlap {
			xEnd := min(xStart + chunkWidth, width)
			yEnd := min(yStart + chunkHeight, height)
	
			wg.Add(1)
			go processChunk(xStart, xEnd, yStart, yEnd)
		}
	}
	


	wg.Wait()
	return SmoothImage(output)
}

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

func SmoothImage(img *image.RGBA) *image.RGBA {
    bounds := img.Bounds()
    smoothed := image.NewRGBA(bounds)
    width, height := bounds.Dx(), bounds.Dy()

    kernel := [][]float64{
        {1, 2, 1},
        {2, 4, 2},
        {1, 2, 1},
    }
    kernelSum := 16.0

    for y := 1; y < height-1; y++ {
        for x := 1; x < width-1; x++ {
            var sumR, sumG, sumB float64
            for ky := -1; ky <= 1; ky++ {
                for kx := -1; kx <= 1; kx++ {
                    nx, ny := x+kx, y+ky
                    r, g, b, _ := img.At(nx, ny).RGBA()
                    weight := kernel[ky+1][kx+1]
                    sumR += float64(r>>8) * weight
                    sumG += float64(g>>8) * weight
                    sumB += float64(b>>8) * weight
                }
            }
            smoothed.Set(x, y, color.RGBA{
                R: uint8(sumR / kernelSum),
                G: uint8(sumG / kernelSum),
                B: uint8(sumB / kernelSum),
                A: 255,
            })
        }
    }
    return smoothed
}
