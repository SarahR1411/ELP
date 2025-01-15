package restoration

import (
    "image"
    "image/color"
    "math"
	"os"
	"log"
	"image/jpeg"
)

func CreateMask(img image.Image, outputPath string) ([][]float64, error) {
	// Get the dimensions of the image
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// create a new mask (2D array of floats) where 1.0 represents a 'damaged' pixel and 0.0 a normal one
	mask := make([][]float64, height)
	for y := 0; y < height; y++ {
		mask[y] = make([]float64, width)
	}

	// Create a new image for the mask
	maskImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// Process each pixel to create the mask
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			originalPixel := img.At(x, y)
			r, g, b, _ := originalPixel.RGBA()

			// convert pixel values to 8-bit range
			r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

			sum := uint32(r8) + uint32(g8) + uint32(b8) // have to use uint32 or else you get an overflow error when comparing with values higher than 255

			// detect white pixels
			if sum > 427 { // to be adjusted
				mask[y][x] = 1.0 // Damaged pixel
				maskImg.Set(x, y, color.White)
			} else {
				mask[y][x] = 0.0 // Normal pixel
				maskImg.Set(x, y, color.Black)
			}
		}
	}

	// Save the mask image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	if err := jpeg.Encode(outputFile, maskImg, nil); err != nil {
		return nil, err
	}

	log.Printf("Mask image saved as %s\n", outputPath)
	return mask, nil
}

func FeatherMask(mask [][]float64, radius int, edgeMask [][]float64) [][]float64 {
    height := len(mask)
	width := len(mask[0])
	feathered := make([][]float64, height)
	for y := range feathered {
		feathered[y] = make([]float64, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if mask[y][x] > 0 {
				feathered[y][x] = 1.0
				for dy := -radius; dy <= radius; dy++ {
					for dx := -radius; dx <= radius; dx++ {
						nx, ny := x+dx, y+dy
						if nx >= 0 && nx < width && ny >= 0 && ny < height {
							distance := float64(dx*dx + dy*dy)
							weight := math.Exp(-distance / float64(radius*radius)) * (1.0 - edgeMask[ny][nx])
							feathered[ny][nx] = math.Max(feathered[ny][nx], weight)
						}
					}
				}
			}
		}
	}
	return feathered
}

func SaveFeatheredMask(feathered [][]float64, outputPath string) error {
    height := len(feathered)
    width := len(feathered[0])
    featheredImg := image.NewGray(image.Rect(0, 0, width, height))
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            intensity := uint8(feathered[y][x] * 255) // Normalize to [0, 255]
            featheredImg.SetGray(x, y, color.Gray{Y: intensity})
        }
    }
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()
    return jpeg.Encode(file, featheredImg, nil)
}



