// $GOPATH /Users/emmafarigoule/go/photo_restoration
//imagePath := "/Users/emmafarigoule/Desktop/old_photo.jpg"

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"sort"
	"time"
)

// Read an image from a file
func loadImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Save an image to a file
func saveImage(img image.Image, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Save the image as JPEG
	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return err
	}
	return nil
}

// Detect scratches and stains using thresholding
func CreateMask(img image.Image, outputPath string) error {
	// Get the dimensions of the image
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create a new image for the mask
	mask := image.NewRGBA(image.Rect(0, 0, width, height))

	// Process each pixel to create the mask
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			originalPixel := img.At(x, y)
			r, g, b, _ := originalPixel.RGBA()

			// Convert pixel values to 8-bit range
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Check if the sum of RGB values is greater than 0
			if r8+g8+b8 > 0 {
				mask.Set(x, y, color.RGBA{0, 0, 0, 255}) // Set to black
			} else {
				mask.Set(x, y, color.RGBA{255, 255, 255, 255}) // Set to white
			}
		}
	}
	// Save the mask image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if err := jpeg.Encode(outputFile, mask, nil); err != nil {
		return err
	}

	log.Printf("Mask image saved as %s\n", outputPath)
	return nil
}

// We'll use Histogram equalization for color correction
func HistEqual(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImg := image.NewRGBA(bounds)

	// Create histograms for each channel
	histR, histG, histB := make([]int, 256), make([]int, 256), make([]int, 256)

	// Calculate the histogram for each channel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			histR[r>>8]++
			histG[g>>8]++
			histB[b>>8]++
		}
	}

	// Calculate the cumulative distribution function (CDF)
	cdfR := computeCDF(histR)
	cdfG := computeCDF(histG)
	cdfB := computeCDF(histB)

	// Find min and max CDF values for normalization
	minR, maxR := findMinMax(cdfR)
	minG, maxG := findMinMax(cdfG)
	minB, maxB := findMinMax(cdfB)

	// Apply histogram equalization to each pixel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			newR := uint8(((cdfR[r>>8] - minR) * 255) / (maxR - minR))
			newG := uint8(((cdfG[g>>8] - minG) * 255) / (maxG - minG))
			newB := uint8(((cdfB[b>>8] - minB) * 255) / (maxB - minB))
			newImg.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)})
		}
	}

	return newImg
}

// computeCDF calculates the cumulative distribution function from a histogram
func computeCDF(hist []int) []int {
	cdf := make([]int, len(hist))
	cdf[0] = hist[0]
	for i := 1; i < len(hist); i++ {
		cdf[i] = cdf[i-1] + hist[i]
	}
	return cdf
}

// findMinMax finds the minimum and maximum non-zero values in a CDF
// Helps avoid zero division error and improves contrast
func findMinMax(cdf []int) (min, max int) {
	min, max = -1, -1
	for _, value := range cdf {
		if value > 0 && min == -1 {
			min = value
		}
		if value > 0 {
			max = value
		}
	}
	return min, max
}

// Get the median color of surrounding pixels for in painting
func getMedianColor(img image.Image, x, y int) color.Color {
	var rValues, gValues, bValues []uint32
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Collect the surrounding pixels' colors in a 3x3 neighborhood
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				c := img.At(nx, ny)
				r, g, b, _ := c.RGBA()
				rValues = append(rValues, r)
				gValues = append(gValues, g)
				bValues = append(bValues, b)
			}
		}
	}

	// Sort the color values to find the median
	sort.Slice(rValues, func(i, j int) bool { return rValues[i] < rValues[j] })
	sort.Slice(gValues, func(i, j int) bool { return gValues[i] < gValues[j] })
	sort.Slice(bValues, func(i, j int) bool { return bValues[i] < bValues[j] })

	// Get the median values
	medianR := uint8(rValues[len(rValues)/2] >> 8)
	medianG := uint8(gValues[len(gValues)/2] >> 8)
	medianB := uint8(bValues[len(bValues)/2] >> 8)

	return color.RGBA{R: medianR, G: medianG, B: medianB, A: 255}
}

// Apply smoothing or filtering to reduce blur
func applySmoothing(img image.Image) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImg := image.NewRGBA(bounds)

	// Simple Gaussian smoothing by averaging neighboring pixels
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Simple average of surrounding pixels
			restoredColor := getMedianColor(img, x, y)
			newImg.Set(x, y, restoredColor)
		}
	}

	return newImg
}

/*
// Inpaint colors the white lines and stains in an image with the median color of surrounding pixels.
func inpaint(img image.Image, mask [][]bool) *image.RGBA {
    bounds := img.Bounds()
    output := image.NewRGBA(bounds)

    // Copy the original image to the output image
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            output.Set(x, y, img.At(x, y))
        }
    }

    // Replace white pixels with the median color of their neighbors
    for y, row := range mask {
        for x, isWhite := range row {
            if isWhite {
                medianColor := getMedianColor(img, x, y)
                output.Set(x, y, medianColor)
            }
        }
    }

    return output
}
*/

// main is the entry point of the application. It performs the following steps:
func main() {
	imagePath := "/Users/emmafarigoule/Desktop/old_photo.jpg"

	restoredImagePath := "restored_photo.jpg"

	// 1. Load the image
	img, err := loadImage(imagePath)
	if err != nil {
		log.Fatalf("Error loading image: %v\n", err)
	}

	start := time.Now()

	// 2. Detect white lines and stains (create a mask)
	// Call the CreateMask function
	outputPath := "new_photo_mask.jpeg"

	err = CreateMask(img, outputPath)
	if err != nil {
		log.Fatalf("Error creating mask: %v", err)
	}

	// 3. Apply scratch removal filter aka inpaint
	//img = inpaint(mask)

	// 4. Apply color correction filter
	color_corrected := HistEqual(img)

	// 5. Apply smoothing to reduce blur
	restoredImg := applySmoothing(color_corrected)

	elapsed := time.Since(start)

	// 6. Save the restored image
	err = saveImage(restoredImg, restoredImagePath)
	if err != nil {
		log.Fatalf("Error saving image: %v\n", err)
	}

	fmt.Printf("Restored image saved to: %s\n", restoredImagePath)
	fmt.Printf("Image restoration took: %v\n", elapsed)
}
