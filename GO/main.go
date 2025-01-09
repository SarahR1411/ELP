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
func detectWhiteLinesAndStains(img image.Image) [][]bool {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	mask := make([][]bool, height)

	for y := 0; y < height; y++ {
		mask[y] = make([]bool, width)
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			// Convert RGBA to grayscale (average of RGB values)
			gray := uint8((r + g + b) / 3)

			// Detect high contrast areas that could represent folds or stains
			if gray > 200 { // These are likely fold lines or light stains
				mask[y][x] = true
			}
		}
	}
	return mask
}

// Inpaint the detected areas using surrounding pixels' colors
/*
func restoreImage(img image.Image, mask [][]bool) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImg := image.NewRGBA(bounds)

	// Iterate over every pixel in the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// If it's part of the scratch or stain (based on the mask), apply restoration
			if mask[y][x] {
				// Apply median filter or other methods for restoration
				restoredColor := getMedianColor(img, x, y)
				newImg.Set(x, y, restoredColor)
			} else {
				// If not part of a scratch, keep the original pixel
				newImg.Set(x, y, img.At(x, y))
			}
		}
	}
	return newImg
}*/

// We'll use Histogram equalization for color correction 

func HistEqual(img image.Image) *image.RGBA {
    bounds := img.Bounds()
    width, height := bounds.Max.X, bounds.Max.Y
    newImg := image.NewRGBA(bounds)

    // Create separate histograms for each channel
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

    // Calculate the cumulative distribution function for each channel
    cdfR := computeCDF(histR)
    cdfG := computeCDF(histG)
    cdfB := computeCDF(histB)

    // Apply histogram equalization to each pixel
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            newR := uint8((cdfR[r>>8] * 255) / cdfR[len(cdfR)-1])
            newG := uint8((cdfG[g>>8] * 255) / cdfG[len(cdfG)-1])
            newB := uint8((cdfB[b>>8] * 255) / cdfB[len(cdfB)-1])
            newImg.SetRGBA(x, y, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)})
        }
    }

    return newImg
}


// Function to calculate the cumulative distribution function

func computeCDF(hist []int) []int {
	cdf := make([]int, len(hist))
	cdf[0] = hist[0]
	for i := 1; i < len(hist); i++ {
		cdf[i] = cdf[i-1] + hist[i]
	}
	return cdf
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

func main() {
	imagePath := "/Users/emmafarigoule/Desktop/old_photo.jpg"

	restoredImagePath := "restored_photo.jpg"

	// Load the image
	img, err := loadImage(imagePath)
	if err != nil {
		log.Fatalf("Error loading image: %v\n", err)
	}

	// Detect white lines and stains (create a mask)
	mask := detectWhiteLinesAndStains(img)

	// Restore the image (remove white lines and stains)
	restoredImg := HistEqual(img) //

	// Apply smoothing to reduce blur
	restoredImg = applySmoothing(restoredImg)

	// Save the restored image
	err = saveImage(restoredImg, restoredImagePath)
	if err != nil {
		log.Fatalf("Error saving image: %v\n", err)
	}

	fmt.Printf("Restored image saved to: %s\n", restoredImagePath)
}
