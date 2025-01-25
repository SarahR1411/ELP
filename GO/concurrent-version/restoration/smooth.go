package restoration

import (
	"image"
	"image/color"
	"math"
	"sync"
)

// Apply gaussian blur and sharpening

func ApplySmoothing(img image.Image) image.Image {

	kernelSize := 3 // smaller kernel = finer smoothing
	sigma := 0.5    // medium smoothing
	blurredImage := GaussianBlur(img, kernelSize, sigma)

	sharpenedImage := Sharpen(blurredImage)

	return sharpenedImage

}

func PostProcessSharpenByChunks(img image.Image, numWorkers int) image.Image {
    bounds := img.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    output := image.NewRGBA(bounds)

    // Sharpen kernel
    kernel := [][]float64{
        {0, -1, 0},
        {-1, 5, -1},
        {0, -1, 0},
    }

    offset := len(kernel) / 2 // Kernel size offset

    // Calculate chunk dimensions with overlap
    chunkWidth := int(math.Ceil(float64(width) / math.Sqrt(float64(numWorkers))))
    chunkHeight := int(math.Ceil(float64(height) / math.Sqrt(float64(numWorkers))))
    overlap := offset // Ensure overlap equals kernel offset

    var wg sync.WaitGroup

    // Worker function for processing a chunk
    processChunk := func(xStart, xEnd, yStart, yEnd int) {
        defer wg.Done()
        for y := max(yStart, offset); y < min(yEnd, height-offset); y++ {
            for x := max(xStart, offset); x < min(xEnd, width-offset); x++ {
                var r, g, b float64
                for ky := -offset; ky <= offset; ky++ {
                    for kx := -offset; kx <= offset; kx++ {
                        nx, ny := x+kx, y+ky
                        px := img.At(nx, ny)
                        pr, pg, pb, _ := px.RGBA()
                        weight := kernel[ky+offset][kx+offset]
                        r += float64(pr) * weight
                        g += float64(pg) * weight
                        b += float64(pb) * weight
                    }
                }

                // Helper function to clamp color values to [0, 255]
                clamp := func(value float64) uint8 {
                    if value < 0 {
                        return 0
                    } else if value > 255*256 {
                        return 255
                    }
                    return uint8(value / 256)
                }

                // Set the processed pixel in the output image
                output.Set(x, y, color.RGBA{
                    R: clamp(r),
                    G: clamp(g),
                    B: clamp(b),
                    A: 255,
                })
            }
        }
    }

    // Launch workers for each chunk
    for yStart := 0; yStart < height; yStart += chunkHeight - overlap {
        for xStart := 0; xStart < width; xStart += chunkWidth - overlap {
            xEnd := min(xStart+chunkWidth, width)
            yEnd := min(yStart+chunkHeight, height)

            wg.Add(1)
            go processChunk(max(0, xStart-overlap), min(width, xEnd+overlap), max(0, yStart-overlap), min(height, yEnd+overlap))
        }
    }

    wg.Wait()
    return output
}


// Gaussian blurr for smoothing and then image sharpening

// Generate a Gaussian kernel dynamically
func generateGaussianKernel(size int, sigma float64) [][]float64 {
	kernel := make([][]float64, size)
	for i := range kernel {
		kernel[i] = make([]float64, size)
	}

	center := size / 2
	sum := 0.0 // To normalize the kernel

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			exponent := -((math.Pow(float64(x-center), 2) + math.Pow(float64(y-center), 2)) / (2 * sigma * sigma))
			kernel[y][x] = math.Exp(exponent) / (2 * math.Pi * sigma * sigma)
			sum += kernel[y][x]
		}
	}

	// Normalize the kernel
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kernel[y][x] /= sum
		}
	}

	return kernel
}

// Apply Gaussian blur with a dynamic kernel size
func GaussianBlur(img image.Image, kernelSize int, sigma float64) image.Image {
	if kernelSize%2 == 0 {
		panic("Kernel size must be an odd number")
	}

	kernel := generateGaussianKernel(kernelSize, sigma)
	bounds := img.Bounds()
	// width, height := bounds.Dx(), bounds.Dy()
	blurredImage := image.NewRGBA(bounds)
	offset := kernelSize / 2

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b float64
			for ky := -offset; ky <= offset; ky++ {
				for kx := -offset; kx <= offset; kx++ {
					nx := x + kx
					ny := y + ky

					// Handle edges by replicating border pixels
					if nx < bounds.Min.X {
						nx = bounds.Min.X
					} else if nx >= bounds.Max.X {
						nx = bounds.Max.X - 1
					}
					if ny < bounds.Min.Y {
						ny = bounds.Min.Y
					} else if ny >= bounds.Max.Y {
						ny = bounds.Max.Y - 1
					}

					px := img.At(nx, ny)
					pr, pg, pb, _ := px.RGBA()
					weight := kernel[ky+offset][kx+offset]
					r += float64(pr) * weight
					g += float64(pg) * weight
					b += float64(pb) * weight
				}
			}
			blurredImage.Set(x, y, color.RGBA{
				R: uint8(r / 256),
				G: uint8(g / 256),
				B: uint8(b / 256),
				A: 255,
			})
		}
	}

	return blurredImage
}

func Sharpen(img image.Image) image.Image {
	bounds := img.Bounds()
	// width, height := bounds.Dx(), bounds.Dy()
	sharpenedImage := image.NewRGBA(bounds)

	// Define a sharpening kernel
	kernel := [][]float64{
		{-1, -2, -1},
		{-2, 13, -2},
		{-1, -2, -1},
	}

	offset := len(kernel) / 2

	for y := bounds.Min.Y + offset; y < bounds.Max.Y-offset; y++ {
		for x := bounds.Min.X + offset; x < bounds.Max.X-offset; x++ {
			var r, g, b float64
			for ky := -offset; ky <= offset; ky++ {
				for kx := -offset; kx <= offset; kx++ {
					nx := x + kx
					ny := y + ky

					// Get the pixel and apply the kernel
					px := img.At(nx, ny)
					pr, pg, pb, _ := px.RGBA()
					weight := kernel[ky+offset][kx+offset]
					r += float64(pr) * weight
					g += float64(pg) * weight
					b += float64(pb) * weight
				}
			}

			// Clamp the values to ensure they fit in [0, 255]
			clamp := func(value float64) uint8 {
				if value < 0 {
					return 0
				} else if value > 255*256 {
					return 255
				}
				return uint8(value / 256)
			}

			sharpenedImage.Set(x, y, color.RGBA{
				R: clamp(r),
				G: clamp(g),
				B: clamp(b),
				A: 255,
			})
		}
	}

	return sharpenedImage
}
