package main

import (
	"fmt"
	"log"
	"time"
	"runtime"
	"GO/concurrent-version/restoration"
)

func main() {
	numWorkers := runtime.NumCPU()
	fmt.Printf("Number of CPU cores: %d\n", numWorkers)
	imagePath := "old_photo.jpeg"
	maskImagePath := "new_photo_mask.jpeg"
	restoredImagePath := "restored_photo.jpg"

	// Load the image
	img, err := restoration.LoadImage(imagePath)
	if err != nil {
		log.Fatalf("Error loading image: %v\n", err)
	}

	start := time.Now()

	// Create the mask in chunks
	mask, err := restoration.CreateMaskByChunks(img, maskImagePath)

	// Generate edge mask for enhanced blending
	edgeMask := restoration.EdgeDetectionConcurrent(img)

	// Feather the mask
	featheredMask := restoration.FeatherMaskConcurrent(mask, 40, edgeMask)

	// Apply scratch removal in chunks
	restoredImg := restoration.InpaintByChunks(img, featheredMask, edgeMask)

	// Post-process for sharpening in chunks
	finalImg := restoration.PostProcessSharpenByChunks(restoredImg, numWorkers)

	elapsed := time.Since(start)

	// Save the final image
	err = restoration.SaveImage(finalImg, restoredImagePath)
	if err != nil {
		log.Fatalf("Error saving restored image: %v\n", err)
	}

	fmt.Printf("Restored image saved to: %s\n", restoredImagePath)
	fmt.Printf("Processing time: %v\n", elapsed)
}
