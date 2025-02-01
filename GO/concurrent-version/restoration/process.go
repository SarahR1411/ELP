package restoration

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// GetBlendedColorWithEdges computes a blended color by averaging nearby pixels weighted by distance and edge strength.
func GetBlendedColorWithEdges(img image.Image, mask [][]float64, edges [][]float64, x, y int) color.Color {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var sumR, sumG, sumB, weightSum float64
	maxRadius := 5
	adjustedRadius := maxRadius
	// Adjust radius for boundary conditions
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

	// Compute weighted average of neighboring pixels
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
	
	// Avoid division by zero
	if weightSum == 0 {
		return img.At(x, y) // Use the original image color directly 
	}

	// Compute final blended color
	return color.RGBA{
		R: uint8((sumR / weightSum) / 256),
		G: uint8((sumG / weightSum) / 256),
		B: uint8((sumB / weightSum) / 256),
		A: 255,
	}
}

// InpaintByChunks performs image inpainting in parallel using chunk processing.
func InpaintByChunks(img image.Image, mask [][]float64, edges [][]float64, numWorkers int) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	output := image.NewRGBA(bounds)

	// Compute chunk size for dividing the image into sections
	chunkHeight := int(math.Ceil(float64(height) / math.Sqrt(2*float64(numWorkers))))
	chunkWidth := int(math.Ceil(float64(width) / math.Sqrt(2*float64(numWorkers))))
	overlap := 5 // Pixels overlapping between chunks for seamless blending
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	// Worker function for processing a chunk of the image
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

// Divide the image into chunks and process each in a goroutine
	for yStart := 0; yStart < height; yStart += chunkHeight - overlap {
		for xStart := 0; xStart < width; xStart += chunkWidth - overlap {
			xEnd := min(xStart + chunkWidth, width)
			yEnd := min(yStart + chunkHeight, height)
	
			wg.Add(1)
			go processChunk(xStart, xEnd, yStart, yEnd)
		}
	}
	


	wg.Wait()
	return SmoothImageConcurrent(output, numWorkers) // Apply final smoothing step
}

// Utility function to return the maximum of two integers.
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// Utility function to return the minimum of two integers.
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

