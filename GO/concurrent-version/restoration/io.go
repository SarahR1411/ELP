package restoration

import (
    "image"
    "image/jpeg"
    "os"
)

// LoadImage loads an image from the specified file path.
// It opens the file, decodes the image, and returns an image.Image object.
func LoadImage(imagePath string) (image.Image, error) {
    file, err := os.Open(imagePath)
    if err != nil {
        return nil, err // Return error if file cannot be opened
    }
    defer file.Close()  // Ensure the file is closed after function execution

    img, _, err := image.Decode(file)   // Decode the image from file
    return img, err
}

// SaveImage saves an image to the specified file path in JPEG format.
// It creates a new file, encodes the image into JPEG format, and writes it to disk.
func SaveImage(img image.Image, outputPath string) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return err  // Return error if file cannot be created
    }
    defer file.Close()  // Ensure the file is closed after function execution

    return jpeg.Encode(file, img, nil)  // Encode and save the image in JPEG format
}
