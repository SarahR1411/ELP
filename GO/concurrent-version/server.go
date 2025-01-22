package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"GO/concurrent-version/restoration"
)

const port = ":8080" // Server port

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected!")

	// 1. Receive the image size
	var imgSize int64
	err := binary.Read(conn, binary.LittleEndian, &imgSize)
	if err != nil {
		log.Println("Error reading image size:", err)
		return
	}
	fmt.Println("Image size received:", imgSize)

	// 2. Receive the image data
	imgData := make([]byte, imgSize)
	_, err = io.ReadFull(conn, imgData)
	if err != nil {
		log.Println("Error reading image data:", err)
		return
	}
	fmt.Println("Image received!")

	// 3. Save the received image to a temporary file
	tempInput := "temp_input.jpg"
	err = os.WriteFile(tempInput, imgData, 0644)
	if err != nil {
		log.Println("Error saving received image:", err)
		return
	}

	// 4. Process the image using the restoration logic
	fmt.Println("Processing image...")
	img, err := restoration.LoadImage(tempInput)
	if err != nil {
		log.Println("Error loading image:", err)
		return
	}

	maskPath := "temp_mask.jpg"
	mask, err := restoration.CreateMask(img, maskPath)
	if err != nil {
		log.Println("Error creating mask:", err)
		return
	}

	edgeMask := restoration.EdgeDetection(img)
	featheredMask := restoration.FeatherMaskConcurrent(mask, 5, edgeMask)
	restoredImg := restoration.InpaintWithEdges(img, featheredMask, edgeMask)
	finalImg := restoration.PostProcessSharpen(restoredImg)

	outputPath := "temp_output.jpg"
	err = restoration.SaveImage(finalImg, outputPath)
	if err != nil {
		log.Println("Error saving restored image:", err)
		return
	}

	// 5. Read the restored image and send it back to the client
	restoredData, err := os.ReadFile(outputPath)
	if err != nil {
		log.Println("Error reading restored image:", err)
		return
	}

	restoredSize := int64(len(restoredData))
	err = binary.Write(conn, binary.LittleEndian, restoredSize)
	if err != nil {
		log.Println("Error sending restored image size:", err)
		return
	}

	_, err = conn.Write(restoredData)
	if err != nil {
		log.Println("Error sending restored image:", err)
		return
	}

	fmt.Println("Restored image sent to client!")
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
	defer listener.Close()

	fmt.Println("Server is running on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn) // handle connection concurrently
	}
}
