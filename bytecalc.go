package main

import (
	"fmt"
	"math"
)

func calculateDimensions(byteCount int) (width, height int, padding int) {
	// Calculate the number of pixels needed (3 bytes per pixel)
	pixelCount := (byteCount + 2) / 3 // Round up to nearest pixel count

	// Determine width and height
	width = int(math.Ceil(math.Sqrt(float64(pixelCount))))
	height = (pixelCount + width - 1) / width // Round up height

	// Calculate padding bytes needed
	padding = (width * height * 3) - byteCount

	return width, height, padding
}

func main() {
	var byteCount int
	fmt.Print("Enter the number of bytes to store: ")
	_, err := fmt.Scanf("%d", &byteCount)
	if err != nil || byteCount <= 0 {
		fmt.Println("Invalid input. Please enter a positive integer.")
		return
	}

	width, height, padding := calculateDimensions(byteCount)
	fmt.Printf("\nFor %d bytes:\n", byteCount)
	fmt.Printf("Width:  %d\n", width)
	fmt.Printf("Height: %d\n", height)
	fmt.Printf("Padding Bytes: %d\n", padding)
	fmt.Printf("Total Bytes: %d\n", byteCount+padding)
}
