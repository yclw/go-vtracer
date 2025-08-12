package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yclw/go-vtracer"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <input_image> <output_svg>\n", os.Args[0])
		fmt.Println("Example: simple_convert image.jpg output.svg")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// Check if input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Fatalf("Input file does not exist: %s", inputPath)
	}

	fmt.Printf("Converting %s -> %s\n", inputPath, outputPath)

	// Convert using default configuration
	err := vtracer.ConvertFile(inputPath, outputPath, nil)
	if err != nil {
		log.Fatalf("Conversion failed: %v", err)
	}

	// Check output file size
	if info, err := os.Stat(outputPath); err == nil {
		fmt.Printf("Conversion successful! Output file size: %d bytes\n", info.Size())
	} else {
		fmt.Println("Conversion successful!")
	}
}
