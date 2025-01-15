package restoration

import (
    "image"
    "image/jpeg"
    "os"
)

func LoadImage(imagePath string) (image.Image, error) {
    file, err := os.Open(imagePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    return img, err
}

func SaveImage(img image.Image, outputPath string) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    return jpeg.Encode(file, img, nil)
}
