package restoration

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

// Get the median color of surrounding pixels for inpainting
func GetBlendedColorWithEdges(img image.Image, mask [][]float64, edges [][]float64, x, y int) color.Color {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var sumR, sumG, sumB, weightSum float64
	maxRadius := 5
	if x < maxRadius || x >= width-maxRadius || y < maxRadius || y >= height-maxRadius {
		maxRadius = 10 // Increase radius near the edges
	}

	for dy := -maxRadius; dy <= maxRadius; dy++ {
		for dx := -maxRadius; dx <= maxRadius; dx++ {
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < width && ny >= 0 && ny < height && mask[ny][nx] < 1.0 {
				c := img.At(nx, ny)
				r, g, b, _ := c.RGBA()
				edgeWeight := 1.0 - edges[ny][nx]
				distance := float64(dx*dx + dy*dy)
				weight := edgeWeight / (math.Sqrt(distance) + 1e-6) // Use sqrt for closer pixels' stronger influence
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
	chunkHeight := int(math.Ceil(float64(height) / math.Sqrt(2*float64(numWorkers))))
	chunkWidth := int(math.Ceil(float64(width) / math.Sqrt(2*float64(numWorkers))))

	var wg sync.WaitGroup

	// Worker function for processing a chunk of the image
	overlap := int(math.Max(5, float64(chunkHeight)/2)) // Number of pixels to overlap
	processChunk := func(xStart, xEnd, yStart, yEnd int) {
		defer wg.Done()
		fmt.Printf("Processing chunk: xStart=%d, xEnd=%d, yStart=%d, yEnd=%d\n", xStart, xEnd, yStart, yEnd)

		for y := yStart; y < yEnd && y < height; y++ {
			for x := xStart; x < xEnd && x < width; x++ {
				fmt.Printf("Processing pixel: x=%d, y=%d, mask=%f\n", x, y, mask[y][x]) // Log mask value

				if mask[y][x] > 0 {
					output.Set(x, y, GetBlendedColorWithEdges(img, mask, edges, x, y))
				} else {
					output.Set(x, y, img.At(x, y))
				}
			}
		}
	}

	for yStart := 0; yStart < height; yStart += chunkHeight - overlap {
		for xStart := 0; xStart < width; xStart += chunkWidth - overlap {
			xEnd := xStart + chunkWidth
			yEnd := yStart + chunkHeight
			wg.Add(1)
			go processChunk(xStart, xEnd, yStart, yEnd)
		}
	}

	wg.Wait()
	return SmoothImage(output)
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
