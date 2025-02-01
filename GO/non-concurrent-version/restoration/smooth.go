package restoration

import (
	"image"
	"image/color"
	"math"
)

// ApplySmoothing applies both Gaussian blur and sharpening to an image.
// It first blurs the image and then sharpens the result to create a smoother version.
func ApplySmoothing(img image.Image) image.Image {
	// Define kernel size and sigma for Gaussian blur
	kernelSize := 3 // smaller kernel = finer smoothing
	sigma := 0.5    // medium smoothing

	// Apply Gaussian blur
	blurredImage := GaussianBlur(img, kernelSize, sigma)

	// Sharpen the blurred image
	sharpenedImage := Sharpen(blurredImage)

	return sharpenedImage

}

// PostProcessSharpen sharpens the input image using a custom sharpening kernel.
// It processes the image with a 3x3 kernel that enhances edges.
func PostProcessSharpen(img image.Image) image.Image {
    bounds := img.Bounds()
	sharpenedImage := image.NewRGBA(bounds)

	// Sharpening kernel (a common 3x3 kernel)
	kernel := [][]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}

	offset := len(kernel) / 2

	// Process each pixel in the image
	for y := bounds.Min.Y + offset; y < bounds.Max.Y-offset; y++ {
		for x := bounds.Min.X + offset; x < bounds.Max.X-offset; x++ {
			var r, g, b float64
			// Apply the kernel to each surrounding pixel
			for ky := -offset; ky <= offset; ky++ {
				for kx := -offset; kx <= offset; kx++ {
					nx := x + kx
					ny := y + ky

					px := img.At(nx, ny)
					pr, pg, pb, _ := px.RGBA()
					weight := kernel[ky+offset][kx+offset]
					r += float64(pr) * weight
					g += float64(pg) * weight
					b += float64(pb) * weight
				}
			}

			// Clamp pixel values to valid range [0, 255]
			clamp := func(value float64) uint8 {
				if value < 0 {
					return 0
				} else if value > 255*256 {
					return 255
				}
				return uint8(value / 256)
			}

			// Set the new pixel value
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

// generateGaussianKernel generates a 2D Gaussian kernel with a given size and sigma.
// The kernel is used for blurring the image in a way that gives a "smoothing" effect.
func generateGaussianKernel(size int, sigma float64) [][]float64 {
	kernel := make([][]float64, size)
	for i := range kernel {
		kernel[i] = make([]float64, size)
	}

	center := size / 2
	sum := 0.0 // To normalize the kernel

	// Calculate the values for the Gaussian kernel
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Compute the exponent of the Gaussian function
			exponent := -((math.Pow(float64(x-center), 2) + math.Pow(float64(y-center), 2)) / (2 * sigma * sigma))
			kernel[y][x] = math.Exp(exponent) / (2 * math.Pi * sigma * sigma)
			sum += kernel[y][x]
		}
	}

	// Normalize the kernel to ensure the sum of all weights equals 1
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			kernel[y][x] /= sum
		}
	}

	return kernel
}

// GaussianBlur applies a Gaussian blur to the image using a dynamically generated kernel.
// The kernel size and sigma (standard deviation) control the strength and extent of the blur.
func GaussianBlur(img image.Image, kernelSize int, sigma float64) image.Image {
	// Ensure kernel size is odd for symmetry
	if kernelSize%2 == 0 {
		panic("Kernel size must be an odd number")
	}

	// Generate Gaussian kernel
	kernel := generateGaussianKernel(kernelSize, sigma)
	bounds := img.Bounds()
	// width, height := bounds.Dx(), bounds.Dy()
	blurredImage := image.NewRGBA(bounds)
	offset := kernelSize / 2

	// Iterate over every pixel in the image.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b float64
			// Apply the kernel to the current pixel and its neighbors.
			for ky := -offset; ky <= offset; ky++ {
				for kx := -offset; kx <= offset; kx++ {
					nx := x + kx
					ny := y + ky

					// Handle edge pixels by replicating the border pixels (no wraparound).
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

					// Get the pixel color at (nx, ny).
					px := img.At(nx, ny)

					// Convert the pixel color to RGBA values (range [0, 65535] as returned by RGBA()).
					pr, pg, pb, _ := px.RGBA()

					// Get the weight of the current kernel element.
					weight := kernel[ky+offset][kx+offset]

					// Accumulate weighted color values.
					r += float64(pr) * weight
					g += float64(pg) * weight
					b += float64(pb) * weight
				}
			}


			// Clamp the color values to fit in the [0, 255] range and set the pixel in the result image.
			blurredImage.Set(x, y, color.RGBA{
				R: uint8(r / 256),	// Convert back from the [0, 65535] range to [0, 255].
				G: uint8(g / 256),
				B: uint8(b / 256),
				A: 255,	// Full opacity.
			})
		}
	}

	return blurredImage
}

// Sharpen applies a sharpening filter to the input image using a 3x3 kernel.
// The kernel enhances edges by emphasizing differences in neighboring pixels.
func Sharpen(img image.Image) image.Image {
	// Get the bounds of the input image to process each pixel.
	bounds := img.Bounds()

	// Create a new image to store the sharpened result.
	sharpenedImage := image.NewRGBA(bounds)

	// Define a sharpening kernel (a common 3x3 kernel for edge enhancement).
	// This kernel emphasizes the central pixel and subtracts the surrounding ones.
	kernel := [][]float64{
		{-1, -2, -1},
		{-2, 13, -2},
		{-1, -2, -1},
	}

	// Calculate the offset from the center of the kernel.
	offset := len(kernel) / 2

	// Iterate over each pixel in the image, excluding the borders to avoid out-of-bounds access.
	for y := bounds.Min.Y + offset; y < bounds.Max.Y-offset; y++ {
		for x := bounds.Min.X + offset; x < bounds.Max.X-offset; x++ {
			var r, g, b float64

			// Apply the sharpening kernel to the current pixel and its neighbors.
			for ky := -offset; ky <= offset; ky++ {
				for kx := -offset; kx <= offset; kx++ {
					nx := x + kx
					ny := y + ky

					// Get the pixel color at (nx, ny).
					px := img.At(nx, ny)

					// Convert the pixel color to RGBA values.
					pr, pg, pb, _ := px.RGBA()

					// Get the weight of the current kernel element.
					weight := kernel[ky+offset][kx+offset]

					// Accumulate weighted color values.
					r += float64(pr) * weight
					g += float64(pg) * weight
					b += float64(pb) * weight
				}
			}

			// Clamp the color values to ensure they fit in the [0, 255] range.
			clamp := func(value float64) uint8 {
				if value < 0 {
					return 0
				} else if value > 255*256 {
					return 255
				}
				return uint8(value / 256)
			}

			// Set the new sharpened pixel in the result image.
			sharpenedImage.Set(x, y, color.RGBA{
				R: clamp(r),
				G: clamp(g),
				B: clamp(b),
				A: 255,	// Full opacity.
			})
		}
	}

	return sharpenedImage
}