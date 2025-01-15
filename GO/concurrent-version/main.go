package main

import (
	"fmt"
	"log"
	"time"

	"GO/concurrent-version/restoration"
)

func main() {
	imagePath := "old_photo.jpeg"
	maskImagePath := "new_photo_mask.jpeg"
	restoredImagePath := "restored_photo.jpg"

	// Load the image
	img, err := restoration.LoadImage(imagePath)
	if err != nil {
		log.Fatalf("Error loading image: %v\n", err)
	}

	start := time.Now()

	// Create the mask for scratches
	mask, err := restoration.CreateMask(img, maskImagePath)
	if err != nil {
		log.Fatalf("Error creating mask: %v\n", err)
	}

	// Generate edge mask for enhanced blending
	edgeMask := restoration.EdgeDetection(img)

	// Feather the mask
	featheredMask := restoration.FeatherMask(mask, 50, edgeMask)

	err = restoration.SaveFeatheredMask(featheredMask, "feathered_mask.jpg")
	if err != nil {
		log.Fatalf("Error saving feathered mask: %v\n", err)
	}

	// Apply scratch removal with edges
	restoredImg := restoration.InpaintWithEdges(img, featheredMask, edgeMask)

	// Post-process for sharpening
	finalImg := restoration.PostProcessSharpen(restoredImg)

	elapsed := time.Since(start)

	// Save the final image
	err = restoration.SaveImage(finalImg, restoredImagePath)
	if err != nil {
		log.Fatalf("Error saving restored image: %v\n", err)
	}

	fmt.Printf("Restored image saved to: %s\n", restoredImagePath)
	fmt.Printf("Processing time: %v\n", elapsed)
}






