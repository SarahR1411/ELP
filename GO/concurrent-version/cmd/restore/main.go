package main

import (
	"fmt"
	"log"
	"time"
	"runtime"
	"GO/concurrent-version/restoration"
	"os"
	"path/filepath"
)

func main() {
	numWorkers := runtime.NumCPU()
	fmt.Printf("Number of workers used: %d\n", numWorkers)
	rootDir, _ := os.Getwd() // Current directory is cmd/restore
	projectDir := filepath.Dir(filepath.Dir(rootDir)) // Move up two levels to concurrent-version

	imagePath := filepath.Join(projectDir, "assets", "old_photo.jpeg")
	maskImagePath := filepath.Join(projectDir, "assets", "new_photo_mask.jpeg")
	restoredImagePath := filepath.Join(projectDir, "assets", "restored_photo.jpg")

	// Load the image
	img, err := restoration.LoadImage(imagePath)
	if err != nil {
		log.Fatalf("Error loading image: %v\n", err)
	}

	start := time.Now()

	// Create the mask in chunks
	mask, err := restoration.CreateMaskByChunks(img, maskImagePath, numWorkers)
	if err != nil {
		log.Fatalf("Error creating mask: %v\n", err)
	}

	// edge mask for blending
	edgeMask := restoration.EdgeDetectionConcurrent(img, numWorkers)

	// feather the mask
	featheredMask := restoration.FeatherMaskConcurrent(mask, 5, edgeMask, numWorkers)

	// apply scratch removal in chunks
	restoredImg := restoration.InpaintByChunks(img, featheredMask, edgeMask, numWorkers)

	// Apply color correction (histogram equalization)
	colorCorrectedImg := restoration.HistEqualConcurrent(restoredImg, numWorkers)

	// Post-process for sharpening and smoothing
	finalImg := restoration.ApplySmoothing(colorCorrectedImg, numWorkers)

	elapsed := time.Since(start)

	// Save the final image
	err = restoration.SaveImage(finalImg, restoredImagePath)
	if err != nil {
		log.Fatalf("Error saving restored image: %v\n", err)
	}

	fmt.Printf("Restored image saved to: %s\n", restoredImagePath)
	fmt.Printf("Processing time: %v\n", elapsed)
}
