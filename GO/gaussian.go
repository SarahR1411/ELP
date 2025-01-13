package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"strings"
)

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
func gaussianBlur(img image.Image, kernelSize int, sigma float64) image.Image {
	if kernelSize%2 == 0 {
		panic("Kernel size must be an odd number")
	}

	kernel := generateGaussianKernel(kernelSize, sigma)
	bounds := img.Bounds()
	blurredImage := image.NewRGBA(bounds)
	offset := kernelSize / 2

	for y := bounds.Min.Y + offset; y < bounds.Max.Y-offset; y++ {
		for x := bounds.Min.X + offset; x < bounds.Max.X-offset; x++ {
			var r, g, b float64
			for ky := -offset; ky <= offset; ky++ {
				for kx := -offset; kx <= offset; kx++ {
					px := img.At(x+kx, y+ky)
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

func main() {
	// Open the image file
	inputFilePath := "C:\\Users\\K3605\\Desktop\\TC\\EPL2\\Owl.jpg"
	file, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer file.Close()

	// Decode the image
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Apply Gaussian blur
	kernelSize := 15 // Increased kernel size
	sigma := 5.0     // Increased sigma value
	blurredImage := gaussianBlur(img, kernelSize, sigma)

	// Generate output file name
	ext := filepath.Ext(inputFilePath)
	name := strings.TrimSuffix(filepath.Base(inputFilePath), ext)
	outputFilePath := filepath.Join(filepath.Dir(inputFilePath), name+"_blurred"+ext)

	// Save the blurred image
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, blurredImage, nil)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	fmt.Println("Blurred image saved to:", outputFilePath)
}
