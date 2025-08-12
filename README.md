# VTracer Go

[![Go Reference](https://pkg.go.dev/badge/github.com/yclw/go-vtracer.svg)](https://pkg.go.dev/github.com/yclw/go-vtracer)
[![Go Report Card](https://goreportcard.com/badge/github.com/yclw/go-vtracer)](https://goreportcard.com/report/github.com/yclw/go-vtracer)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

This is a **Go language binding** for [VTracer](https://github.com/visioncortex/vtracer), enabling you to easily use VTracer for high-quality image vectorization in your Go applications.

## Features

- **High-Quality Vectorization**: Convert raster images to compact SVG vectors
- **High Performance**: Direct calls to optimized Rust core library
- **Easy to Use**: Idiomatic Go API that's simple and clean
- **Multiple Input Types**: Support for file paths, Go image.Image, and byte arrays
- **Flexible Configuration**: Fully customizable parameters
- **Cross-Platform**: Support for Linux (x86_64/aarch64), macOS (x86_64/arm64), Windows (x86_64)

## Quick Start

### 1. Install Package

```bash
go get github.com/yclw/go-vtracer
```

### 2. Basic Usage

```go
package main

import (
    "log"
    "github.com/yclw/go-vtracer"
)

func main() {
    // Convert with default configuration
    err := vtracer.ConvertFile("input.jpg", "output.svg", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Or use custom configuration
    config := vtracer.DefaultConfig()
    config.ColorPrecision = 8
    config.FilterSpeckle = 8
    err = vtracer.ConvertFile("photo.jpg", "photo.svg", config)
    if err != nil {
        log.Fatal(err)
    }
}
```

## API Reference

### Type Definitions

```go
type ColorMode uint8
const (
    ColorModeColor  ColorMode = 0  // Color mode
    ColorModeBinary ColorMode = 1  // Binary mode
)

type Hierarchical uint8
const (
    HierarchicalStacked Hierarchical = 0  // Stacked mode
    HierarchicalCutout  Hierarchical = 1  // Cutout mode
)

type PathMode uint8
const (
    PathModeNone    PathMode = 0  // Pixel mode
    PathModePolygon PathMode = 1  // Polygon mode
    PathModeSpline  PathMode = 2  // Spline mode
)
```

### Configuration Options

```go
type Config struct {
    ColorMode       ColorMode    // Color processing mode
    Hierarchical    Hierarchical // Hierarchical structure mode
    FilterSpeckle   int          // Filter speckle size (pixels)
    ColorPrecision  int          // Color precision (1-8 bits)
    LayerDifference int          // Layer color difference (0-255)
    Mode            PathMode     // Path fitting mode
    CornerThreshold int          // Corner threshold (0-180 degrees)
    LengthThreshold float64      // Length threshold (3.5-10.0)
    MaxIterations   int          // Maximum iterations
    SpliceThreshold int          // Splice threshold (0-180 degrees)
    PathPrecision   int          // Path precision (decimal places)
}
```

### Main Functions

```go
// Create default configuration
func DefaultConfig() *Config

// Convert image file to SVG file
func ConvertFile(inputPath, outputPath string, config *Config) error

// Convert Go image.Image to SVG string
func ConvertImage(img image.Image, config *Config) (string, error)

// Convert RGBA byte array to SVG string
func ConvertBytes(pixels []byte, width, height int, config *Config) (string, error)
```

## License

MIT

## Contributing

Contributions of all kinds are welcome, such as bug fixes, new features, documentation improvements, etc.

## Credits

- [VTracer](https://github.com/visioncortex/vtracer)
