package main

import (
	"fmt"
	"log"
	"time"

	"GO/restoration"
)

// main is the entry point of the application. It performs the following steps:
// func main() {
// 	imagePath := "old_photo.jpeg"
// 	restoredImagePath := "restored_photo.jpg"
// 	maskImagePath := "new_photo_mask.jpeg"

// 	// 1. Load the image
// 	img, err := loadImage(imagePath)
// 	if err != nil {
// 		log.Fatalf("Error loading image: %v\n", err)
// 	}

// 	start := time.Now()

// 	// 2. Create a mask for scratches and stains
// 	mask, err := CreateMask(img, maskImagePath)
// 	if err != nil {
// 		log.Fatalf("Error creating mask: %v\n", err)
// 	}

// 	// 3. Apply scratch removal (inpainting)
// 	inpaintedImg := inpaint(img, mask)

// 	// 4. Apply color correction
// 	colorCorrectedImg := HistEqual(inpaintedImg)

// 	// 5. Apply smoothing to reduce blur
// 	restoredImg := applySmoothing(colorCorrectedImg)

// 	elapsed := time.Since(start)

// 	// 6. Save the restored image
// 	err = saveImage(restoredImg, restoredImagePath)
// 	if err != nil {
// 		log.Fatalf("Error saving restored image: %v\n", err)
// 	}

// 	fmt.Printf("Restored image saved to: %s\n", restoredImagePath)
// 	fmt.Printf("Image restoration took: %v\n", elapsed)
// }

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






