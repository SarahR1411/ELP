package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const serverAddress = "localhost:8080" // server address

func main() {
	// 1. Connect to the server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error connecting to server: %v\n", err)
	}
	defer conn.Close()

	// 2. Open the image file to send
	imagePath := "old_photo.jpeg"
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("Error reading image file: %v\n", err)
	}

	// 3. Send the image size
	imgSize := int64(len(imageData))
	err = binary.Write(conn, binary.LittleEndian, imgSize)
	if err != nil {
		log.Fatalf("Error sending image size: %v\n", err)
	}

	// 4. Send the image data
	_, err = conn.Write(imageData)
	if err != nil {
		log.Fatalf("Error sending image data: %v\n", err)
	}
	fmt.Println("Image sent to server!")

	// 5. Receive the restored image size
	var restoredSize int64
	err = binary.Read(conn, binary.LittleEndian, &restoredSize)
	if err != nil {
		log.Fatalf("Error reading restored image size: %v\n", err)
	}
	fmt.Println("Restored image size received:", restoredSize)

	// 6. Receive the restored image data
	restoredData := make([]byte, restoredSize)
	_, err = io.ReadFull(conn, restoredData)
	if err != nil {
		log.Fatalf("Error reading restored image data: %v\n", err)
	}

	// 7. Save the restored image to a file
	outputPath := "restored_by_server.jpg"
	err = os.WriteFile(outputPath, restoredData, 0644)
	if err != nil {
		log.Fatalf("Error saving restored image: %v\n", err)
	}

	fmt.Printf("Restored image saved to: %s\n", outputPath)
}
