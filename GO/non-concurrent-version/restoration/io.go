package restoration

import (
    "image"
    "image/jpeg"
    "os"
)

// LoadImage loads an image from a specified file path.
// It opens the file, decodes the image, and returns the decoded image.
func LoadImage(imagePath string) (image.Image, error) {
    // Open the image file
    file, err := os.Open(imagePath)
    if err != nil {
        return nil, err // Return an error if the file cannot be opened
    }
    defer file.Close()  // Ensure the file is closed after the function returns

    // Decode the image from the file
    img, _, err := image.Decode(file)
    return img, err // Return the decoded image or an error if decoding fails
}

// SaveImage saves an image to a specified file path in JPEG format.
// It creates a new file and encodes the image into the file.
func SaveImage(img image.Image, outputPath string) error {
    // Create the output file for saving the image
    file, err := os.Create(outputPath)
    if err != nil {
        return err  // Return an error if the file cannot be created
    }
    defer file.Close()  // Ensure the file is closed after the function returns

    // Encode the image into the file in JPEG format
    return jpeg.Encode(file, img, nil)  // Returns an error if encoding fails
}
