package restoration

import (
    "image"
    "image/color"
    "math"
	"os"
	"log"
	"image/jpeg"
)

// CreateMask generates a mask image that marks damaged pixels (e.g., white pixels) as 1.0 and normal pixels as 0.0.
// The mask is returned as a 2D float64 array and the mask image is saved to the specified output path.
func CreateMask(img image.Image, outputPath string) ([][]float64, error) {
	// Get the dimensions of the input image
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Initialize a 2D array for the mask (damaged pixels will be 1.0, normal pixels will be 0.0)
	mask := make([][]float64, height)
	for y := 0; y < height; y++ {
		mask[y] = make([]float64, width)
	}

	// Create a new RGBA image for the mask (for visualization)
	maskImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Process each pixel to create the mask
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			originalPixel := img.At(x, y)
			r, g, b, _ := originalPixel.RGBA()

			// Convert the pixel values to the 8-bit range (0-255)
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			// Sum the RGB values to detect "white" pixels
			sum := uint32(r8) + uint32(g8) + uint32(b8) 

			// Detect if the pixel is white (damaged) based on the sum of RGB values
			if sum > 427 { 	// Threshold for detecting white pixels, can be adjusted
				mask[y][x] = 1.0 // Damaged pixel
				maskImg.Set(x, y, color.White)	// Set mask image pixel to white
			} else {
				mask[y][x] = 0.0 // Normal pixel
				maskImg.Set(x, y, color.Black)	// Set mask image pixel to black
			}
		}
	}

	// Save the generated mask image to the specified file path
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err	// Return error if file cannot be created
	}
	defer outputFile.Close()	// Ensure the file is closed after writing

	// Encode and save the mask image in JPEG format
	if err := jpeg.Encode(outputFile, maskImg, nil); err != nil {
		return nil, err	// Return error if encoding fails
	}

	log.Printf("Mask image saved as %s\n", outputPath)
	return mask, nil	// Return the mask as a 2D float64 array
}

// FeatherMask applies a feathering effect to a binary mask, creating a soft transition between damaged and normal areas.
// The feathering effect is applied by calculating a weighted average of neighboring pixels using a Gaussian function.
func FeatherMask(mask [][]float64, radius int, edgeMask [][]float64) [][]float64 {
    height := len(mask)
	width := len(mask[0])
	feathered := make([][]float64, height)

	// Initialize a new 2D slice for the feathered mask
	for y := range feathered {
		feathered[y] = make([]float64, width)
	}

	// Apply feathering to the mask pixels
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if mask[y][x] > 0 {	// If it's a damaged pixel
				feathered[y][x] = 1.0	// Start with full feathering at the damaged pixel

				// Apply a weighted average to neighboring pixels based on distance and edgeMask
				for dy := -radius; dy <= radius; dy++ {
					for dx := -radius; dx <= radius; dx++ {
						nx, ny := x+dx, y+dy
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							// Calculate the distance and weight using a Gaussian function
							distance := float64(dx*dx + dy*dy)
							weight := math.Exp(-distance / float64(radius*radius)) * (1.0 - edgeMask[ny][nx])
							// Update the feathered value if the new weight is higher
							feathered[ny][nx] = math.Max(feathered[ny][nx], weight)
						}
					}
				}
			}
		}
	}
	// Return the feathered mask
	return feathered
}

// SaveFeatheredMask saves the feathered mask as an image file (grayscale).
// The feathered mask is converted to an 8-bit grayscale image and then saved as a JPEG.
func SaveFeatheredMask(feathered [][]float64, outputPath string) error {
    height := len(feathered)
    width := len(feathered[0])
    featheredImg := image.NewGray(image.Rect(0, 0, width, height))
	// Convert the feathered mask to grayscale image
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            intensity := uint8(feathered[y][x] * 255) // Normalize the feathered value to [0, 255]
            featheredImg.SetGray(x, y, color.Gray{Y: intensity})
        }
    }

	// Create and open the output file for saving the image
    file, err := os.Create(outputPath)
    if err != nil {
        return err	// Return error if file cannot be created
    }
    defer file.Close()	// Ensure the file is closed after writing

	// Encode and save the feathered mask image in JPEG format
    return jpeg.Encode(file, featheredImg, nil)
}



