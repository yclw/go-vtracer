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

## Examples

### Convert File with Custom Settings

```go
package main

import (
    "log"
    "github.com/yclw/go-vtracer"
)

func main() {
    config := vtracer.DefaultConfig()
    
    // High quality settings
    config.ColorPrecision = 8      // Maximum color precision
    config.FilterSpeckle = 4       // Remove small noise
    config.Mode = vtracer.PathModeSpline  // Use spline curves
    config.CornerThreshold = 60    // Detect corners at 60 degrees
    config.LengthThreshold = 4.0   // Minimum curve length
    
    err := vtracer.ConvertFile("input.png", "output.svg", config)
    if err != nil {
        log.Fatal("Conversion failed:", err)
    }
    
    log.Println("Successfully converted to SVG!")
}
```

### Convert Go Image

```go
package main

import (
    "image"
    "image/jpeg"
    "log"
    "os"
    "github.com/yclw/go-vtracer"
)

func main() {
    // Open and decode image
    file, err := os.Open("input.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    
    img, err := jpeg.Decode(file)
    if err != nil {
        log.Fatal(err)
    }
    
    // Convert to SVG
    svgContent, err := vtracer.ConvertImage(img, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Save SVG
    err = os.WriteFile("output.svg", []byte(svgContent), 0644)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Image converted and saved as SVG!")
}
```

### Binary Image Processing

```go
config := vtracer.DefaultConfig()
config.ColorMode = vtracer.ColorModeBinary  // Binary mode for logos/line art
config.FilterSpeckle = 16                   // Remove larger noise spots

err := vtracer.ConvertFile("logo.png", "logo.svg", config)
if err != nil {
    log.Fatal(err)
}
```

## Platform Support

The library includes pre-built native libraries for:

- **Linux**: x86_64, aarch64
- **macOS**: x86_64 (Intel), arm64 (Apple Silicon)  
- **Windows**: x86_64

The appropriate library is automatically selected based on your platform.

## Performance Tips

1. **Use Binary Mode**: For logos and line art, use `ColorModeBinary` for better performance
2. **Adjust Color Precision**: Lower values (4-6) are faster, higher values (7-8) are more accurate
3. **Filter Speckles**: Use `FilterSpeckle` to remove noise and reduce output size
4. **Choose Path Mode**: `PathModePolygon` is fastest, `PathModeSpline` gives smoothest results

## Error Handling

The library provides detailed error messages for common issues:

```go
err := vtracer.ConvertFile("input.jpg", "output.svg", nil)
if err != nil {
    log.Printf("Conversion failed: %v", err)
    // Handle specific error cases
}
```

## Building from Source

If you need to build the native library from source:

```bash
# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Build the library
cargo build --release

# Copy to lib directory
cp target/release/libvtracer_go.* lib/
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Credits

- [VTracer](https://github.com/visioncortex/vtracer) - The core vectorization library
- [VisionCortex](https://visioncortex.org/) - Computer vision research group
