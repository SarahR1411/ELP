package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"time"
	"GO/concurrent-version/restoration"
)

const port = ":8080" // Server port

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected!")
	numWorkers := runtime.NumCPU()
	fmt.Printf("Number of workers used: %d\n", numWorkers)

	// Start timing the processing
	start := time.Now()

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

	// Generate unique file names for this connection
	timestamp := time.Now().UnixNano()
	tempInput := fmt.Sprintf("temp_input_%d.jpg", timestamp)
	tempMask := fmt.Sprintf("temp_mask_%d.jpg", timestamp)
	tempOutput := fmt.Sprintf("temp_output_%d.jpg", timestamp)

	// 3. Save the received image to a temporary file
	err = os.WriteFile(tempInput, imgData, 0644)
	if err != nil {
		log.Println("Error saving received image:", err)
		return
	}
	defer os.Remove(tempInput) // Clean up the temporary input file

	// 4. Process the image using the restoration logic
	fmt.Println("Processing image...")
	img, err := restoration.LoadImage(tempInput)
	if err != nil {
		log.Println("Error loading image:", err)
		return
	}

	// Create a mask
	mask, err := restoration.CreateMaskByChunks(img, tempMask, numWorkers)
	if err != nil {
		log.Println("Error creating mask:", err)
		return
	}
	defer os.Remove(tempMask) // Clean up the temporary mask file

	// Generate edge mask
	edgeMask := restoration.EdgeDetectionConcurrent(img, numWorkers)

	// Feather the mask
	featheredMask := restoration.FeatherMaskConcurrent(mask, 5, edgeMask, numWorkers)

	// Apply scratch removal
	restoredImg := restoration.InpaintByChunks(img, featheredMask, edgeMask, numWorkers)

	// Apply color correction
	colorCorrectedImg := restoration.HistEqualConcurrent(restoredImg, numWorkers)

	// Post-process for sharpening and smoothing
	finalImg := restoration.ApplySmoothing(colorCorrectedImg, numWorkers)

	// Save the final output
	err = restoration.SaveImage(finalImg, tempOutput)
	if err != nil {
		log.Println("Error saving restored image:", err)
		return
	}
	defer os.Remove(tempOutput) // Clean up the temporary output file

	// Calculate processing time
	elapsed := time.Since(start)
	fmt.Printf("Image processing completed in: %v\n", elapsed)

	// 5. Send metadata (number of workers and processing time) to the client
	metadata := fmt.Sprintf("Workers: %d, Processing Time: %v", numWorkers, elapsed)
	metadataSize := int64(len(metadata))
	err = binary.Write(conn, binary.LittleEndian, metadataSize)
	if err != nil {
		log.Println("Error sending metadata size:", err)
		return
	}

	_, err = conn.Write([]byte(metadata))
	if err != nil {
		log.Println("Error sending metadata:", err)
		return
	}

	// 6. Read the restored image and send it back to the client
	restoredData, err := os.ReadFile(tempOutput)
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
		go handleConnection(conn) // Handle each connection concurrently
	}
}
