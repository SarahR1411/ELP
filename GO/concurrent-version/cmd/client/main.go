package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// To use activate the server first by running the server file in a seperate terminal
// Then run the client file using a command like this: go run main.go -file /Users/sarah/Desktop/ELP_S1/ELP/GO/concurrent-version/assets/old_photo.jpeg
// If you have issues geting the correct path write the pwd command in the terminal for help

const serverAddress = "localhost:8080" // server address

func main() {
	// Parse the input file path from the command line
	imagePath := flag.String("file", "", "Path to the image file to send")
	flag.Parse()

	if *imagePath == "" {
		log.Fatal("Please provide the path to the image file using the -file flag")
	}

	// 1. Connect to the server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Error connecting to server: %v\n", err)
	}
	defer conn.Close()

	// 2. Open the image file to send
	imageData, err := os.ReadFile(*imagePath)
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

	// 5. Receive metadata size
	var metadataSize int64
	err = binary.Read(conn, binary.LittleEndian, &metadataSize)
	if err != nil {
		log.Fatalf("Error reading metadata size: %v\n", err)
	}

	// 6. Receive metadata
	metadata := make([]byte, metadataSize)
	_, err = io.ReadFull(conn, metadata)
	if err != nil {
		log.Fatalf("Error reading metadata: %v\n", err)
	}
	fmt.Println("Metadata received:", string(metadata))

	// 7. Receive the restored image size
	var restoredSize int64
	err = binary.Read(conn, binary.LittleEndian, &restoredSize)
	if err != nil {
		log.Fatalf("Error reading restored image size: %v\n", err)
	}
	fmt.Println("Restored image size received:", restoredSize)

	// 8. Receive the restored image data
	restoredData := make([]byte, restoredSize)
	_, err = io.ReadFull(conn, restoredData)
	if err != nil {
		log.Fatalf("Error reading restored image data: %v\n", err)
	}

	// 9. Save the restored image to a file
	outputPath := "restored_by_server.jpg"
	err = os.WriteFile(outputPath, restoredData, 0644)
	if err != nil {
		log.Fatalf("Error saving restored image: %v\n", err)
	}

	fmt.Printf("Restored image saved to: %s\n", outputPath)
}
