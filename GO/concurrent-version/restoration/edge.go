package restoration

import (
	"image"
	"math"
	"sync"
	"runtime"
)



func EdgeDetectionConcurrent(img image.Image) [][]float64 {
    bounds := img.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    edges := make([][]float64, height)
    for i := range edges {
        edges[i] = make([]float64, width)
    }

    // Sobel kernels
    sobelX := [][]int{
        {-1, 0, 1},
        {-2, 0, 2},
        {-1, 0, 1},
    }
    sobelY := [][]int{
        {-1, -2, -1},
        {0, 0, 0},
        {1, 2, 1},
    }

    var maxGradient float64
    maxGradientMutex := &sync.Mutex{} // Protects access to maxGradient

    // Worker function to process rows
    processRow := func(yStart, yEnd int, wg *sync.WaitGroup) {
        defer wg.Done()
        for y := yStart; y < yEnd; y++ {
            for x := 1; x < width-1; x++ {
                var gx, gy float64
                for ky := -1; ky <= 1; ky++ {
                    for kx := -1; kx <= 1; kx++ {
                        px := img.At(x+kx, y+ky)
                        r, g, b, _ := px.RGBA()
                        gray := float64(r+g+b) / (3.0 * 256.0)
                        gx += gray * float64(sobelX[ky+1][kx+1])
                        gy += gray * float64(sobelY[ky+1][kx+1])
                    }
                }
                gradient := math.Sqrt(gx*gx + gy*gy)
                edges[y][x] = gradient

                // Safely update maxGradient
                maxGradientMutex.Lock()
                if gradient > maxGradient {
                    maxGradient = gradient
                }
                maxGradientMutex.Unlock()
            }
        }
    }

    // Divide work among workers
    numWorkers := runtime.NumCPU()
    wg := &sync.WaitGroup{}
    rowsPerWorker := height / numWorkers
    for i := 0; i < numWorkers; i++ {
        yStart := i * rowsPerWorker
        yEnd := (i + 1) * rowsPerWorker
        if i == numWorkers-1 { // Handle remaining rows for the last worker
            yEnd = height
        }
        wg.Add(1)
        go processRow(yStart, yEnd, wg)
    }
    wg.Wait()

    // Normalize and apply threshold
    threshold := 0.2 // do more tests maybe to find good value
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            edges[y][x] /= maxGradient
            if edges[y][x] < threshold {
                edges[y][x] = 0.0
            }
        }
    }

    return edges
}

