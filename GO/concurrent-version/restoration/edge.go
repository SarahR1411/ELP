package restoration

import (
	"image"
	"math"
)



func EdgeDetection(img image.Image) [][]float64 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	edges := make([][]float64, height)
	for i := range edges {
		edges[i] = make([]float64, width)
	}

	// Sobel kernels
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
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var gx, gy float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := img.At(x+kx, y+ky)
					r, g, b, _ := px.RGBA()
					gray := float64(r+g+b) / (3.0 * 256.0)
					gx += gray * float64(sobelX[ky+1][kx+1])
					gy += gray * float64(sobelY[ky+1][kx+1])
				}
			}
			gradient := math.Sqrt(gx*gx + gy*gy)
			if gradient > maxGradient {
				maxGradient = gradient
			}
			edges[y][x] = gradient
		}
	}

	// Normalize and apply threshold
	threshold := 0.2 // Do more tests to find good value
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			edges[y][x] /= maxGradient
			if edges[y][x] < threshold {
				edges[y][x] = 0.0
			}
		}
	}

	return edges
}
