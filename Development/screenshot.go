package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"
)

const (
	PixelSize    = 10 // Size of each pixel
	PixelSpacing = 1  // Spacing between pixels
)

func main() {
	var ipAddress string

	fmt.Print("Enter the IP address of the API endpoint: ")
	fmt.Scan(&ipAddress)

	// Make a GET request to the API endpoint
	response, err := http.Get(fmt.Sprintf("http://%s/api/screen", ipAddress))
	if err != nil {
		fmt.Println("Failed to fetch RGB data from the API:", err)
		return
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("API request returned a non-OK status:", response.StatusCode)
		return
	}

	// Extract the RGB data from the response body
	var rgbData []uint32
	if err := json.NewDecoder(response.Body).Decode(&rgbData); err != nil {
		fmt.Println("Failed to decode RGB data:", err)
		return
	}

	// Calculate the width and height of the image based on the RGB data length
	width := 32
	height := len(rgbData) / width

	// Calculate the image size with spacing
	imageWidth := (PixelSize+PixelSpacing)*width + PixelSpacing
	imageHeight := (PixelSize+PixelSpacing)*height + PixelSpacing

	// Create a new RGBA image
	img := image.NewRGBA(image.Rect(0, 0, int(imageWidth), int(imageHeight)))

	// Fill the image with black color
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	// Loop through the RGB data and set the pixels in the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := y*width + x
			rgb := rgbData[index]

			r := uint8((rgb >> 16) & 0xFF)
			g := uint8((rgb >> 8) & 0xFF)
			b := uint8(rgb & 0xFF)
			a := uint8(0xFF)

			// Calculate the position of the pixel with spacing
			posX := (PixelSize+PixelSpacing)*float64(x) + PixelSpacing
			posY := (PixelSize+PixelSpacing)*float64(y) + PixelSpacing

			// Draw the pixel with the specified color
			drawPixel(img, posX, posY, r, g, b, a)
		}
	}

	// Create the output file
	file, err := os.Create("pixelized_image.png")
	if err != nil {
		fmt.Println("Failed to create output file:", err)
		return
	}
	defer file.Close()

	// Encode the image as PNG and save it to the file
	err = png.Encode(file, img)
	if err != nil {
		fmt.Println("Failed to encode and save the image:", err)
		return
	}

	fmt.Println("Pixelized image created and saved as pixelized_image.png")
}

// Draw a pixel with the specified color at the given position
func drawPixel(img *image.RGBA, x, y float64, r, g, b, a uint8) {
	for i := 0; i < PixelSize; i++ {
		for j := 0; j < PixelSize; j++ {
			// Calculate the pixel position
			pixelX := int(x) + i
			pixelY := int(y) + j

			// Set the color of the pixel
			img.Set(pixelX, pixelY, color.RGBA{r, g, b, a})
		}
	}
}
