package restoration

import (
	"image"
	"math"
)

// EdgeDetection applies the Sobel edge detection algorithm to an image.
// The Sobel operator is used to detect edges by calculating the gradient of pixel intensities.
func EdgeDetection(img image.Image) [][]float64 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Initialize a 2D slice to store edge values (gradient magnitudes) for each pixel
	edges := make([][]float64, height)
	for i := range edges {
		edges[i] = make([]float64, width)
	}

	// Sobel kernels for detecting edges in the X and Y directions
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

	// Loop through each pixel in the image (excluding the borders)
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var gx, gy float64

			// Apply the Sobel kernel to the surrounding 3x3 neighborhood for the current pixel
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := img.At(x+kx, y+ky)
					r, g, b, _ := px.RGBA()

					// Convert the pixel's RGB values to grayscale using an average method
					gray := float64(r+g+b) / (3.0 * 256.0)

					// Calculate the gradients in the X and Y directions using the Sobel kernels
					gx += gray * float64(sobelX[ky+1][kx+1])
					gy += gray * float64(sobelY[ky+1][kx+1])
				}
			}

			// Calculate the gradient magnitude (edge strength) using Pythagoras' theorem
			gradient := math.Sqrt(gx*gx + gy*gy)

			// Track the maximum gradient value to normalize later
			if gradient > maxGradient {
				maxGradient = gradient
			}
			// Store the calculated gradient for the current pixel
			edges[y][x] = gradient
		}
	}

	// Normalize the edge magnitudes to a range of 0 to 1 and apply a threshold
	threshold := 0.2 // This threshold can be adjusted for better results
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Normalize the gradient by dividing by the maximum gradient value
			edges[y][x] /= maxGradient
			// Set values below the threshold to 0 (indicating no significant edge)
			if edges[y][x] < threshold {
				edges[y][x] = 0.0
			}
		}
	}
	
	// Return the 2D array of edge magnitudes after thresholding and normalization
	return edges
}
