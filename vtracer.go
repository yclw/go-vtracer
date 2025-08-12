package vtracer

/*
#cgo linux,amd64 LDFLAGS: ${SRCDIR}/lib/libvtracer_go_linux_amd64.so -Wl,-rpath,${SRCDIR}/lib
#cgo linux,arm64 LDFLAGS: ${SRCDIR}/lib/libvtracer_go_linux_arm64.so -Wl,-rpath,${SRCDIR}/lib
#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/lib/libvtracer_go_darwin_amd64.dylib -Wl,-rpath,${SRCDIR}/lib
#cgo darwin,arm64 LDFLAGS: ${SRCDIR}/lib/libvtracer_go_darwin_arm64.dylib -Wl,-rpath,${SRCDIR}/lib
#cgo windows,amd64 LDFLAGS: ${SRCDIR}/lib/vtracer_go_windows_amd64.dll
#include <stdlib.h>

typedef struct {
    unsigned char color_mode;
    unsigned char hierarchical;
    size_t filter_speckle;
    int color_precision;
    int layer_difference;
    unsigned char mode;
    int corner_threshold;
    double length_threshold;
    size_t max_iterations;
    int splice_threshold;
    unsigned int path_precision;
} VtracerConfig;

int vtracer_convert_file(const char* input_path, const char* output_path, const VtracerConfig* config);
int vtracer_convert_bytes(const unsigned char* image_data, size_t image_len, size_t width, size_t height, const VtracerConfig* config, char** output);
void vtracer_free_string(char* s);
VtracerConfig vtracer_default_config();
*/
import "C"
import (
	"fmt"
	"image"
	"unsafe"
)

// ColorMode represents the color processing mode
type ColorMode uint8

const (
	ColorModeColor  ColorMode = 0 // Color mode
	ColorModeBinary ColorMode = 1 // Binary mode
)

// Hierarchical represents the hierarchical processing mode
type Hierarchical uint8

const (
	HierarchicalStacked Hierarchical = 0 // Stacked mode
	HierarchicalCutout  Hierarchical = 1 // Cutout mode
)

// PathMode represents the path fitting mode
type PathMode uint8

const (
	PathModeNone    PathMode = 0 // Pixel mode
	PathModePolygon PathMode = 1 // Polygon mode
	PathModeSpline  PathMode = 2 // Spline mode
)

// Error code mapping table
var errorMessages = map[int]string{
	-1: "invalid input parameters",
	-2: "invalid input path or pixel data length error",
	-3: "invalid output path or string conversion failed",
	-4: "conversion failed",
}

// codeToError converts C function return codes to Go errors
func codeToError(code int) error {
	if code == 0 {
		return nil
	}
	if msg, exists := errorMessages[code]; exists {
		return fmt.Errorf("vtracer: %s (错误码: %d)", msg, code)
	}
	return fmt.Errorf("vtracer: 未知错误 (错误码: %d)", code)
}

// Config contains vtracer configuration options
type Config struct {
	ColorMode       ColorMode    // Color mode: color or binary
	Hierarchical    Hierarchical // Hierarchical structure: stacked or cutout
	FilterSpeckle   int          // Filter speckle size
	ColorPrecision  int          // Color precision (1-8)
	LayerDifference int          // Layer difference (0-255)
	Mode            PathMode     // Path fitting mode
	CornerThreshold int          // Corner threshold (0-180 degrees)
	LengthThreshold float64      // Length threshold (3.5-10.0)
	MaxIterations   int          // Maximum iterations
	SpliceThreshold int          // Splice threshold (0-180 degrees)
	PathPrecision   int          // Path precision
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	cConfig := C.vtracer_default_config()
	return &Config{
		ColorMode:       ColorMode(cConfig.color_mode),
		Hierarchical:    Hierarchical(cConfig.hierarchical),
		FilterSpeckle:   int(cConfig.filter_speckle),
		ColorPrecision:  int(cConfig.color_precision),
		LayerDifference: int(cConfig.layer_difference),
		Mode:            PathMode(cConfig.mode),
		CornerThreshold: int(cConfig.corner_threshold),
		LengthThreshold: float64(cConfig.length_threshold),
		MaxIterations:   int(cConfig.max_iterations),
		SpliceThreshold: int(cConfig.splice_threshold),
		PathPrecision:   int(cConfig.path_precision),
	}
}

// toCConfig converts Go configuration to C configuration
func (c *Config) toCConfig() C.VtracerConfig {
	return C.VtracerConfig{
		color_mode:       C.uchar(c.ColorMode),
		hierarchical:     C.uchar(c.Hierarchical),
		filter_speckle:   C.size_t(c.FilterSpeckle),
		color_precision:  C.int(c.ColorPrecision),
		layer_difference: C.int(c.LayerDifference),
		mode:             C.uchar(c.Mode),
		corner_threshold: C.int(c.CornerThreshold),
		length_threshold: C.double(c.LengthThreshold),
		max_iterations:   C.size_t(c.MaxIterations),
		splice_threshold: C.int(c.SpliceThreshold),
		path_precision:   C.uint(c.PathPrecision),
	}
}

// ConvertFile converts an image file to SVG file
func ConvertFile(inputPath, outputPath string, config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	cInputPath := C.CString(inputPath)
	cOutputPath := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cInputPath))
	defer C.free(unsafe.Pointer(cOutputPath))

	cConfig := config.toCConfig()
	result := C.vtracer_convert_file(cInputPath, cOutputPath, &cConfig)

	return codeToError(int(result))
}

// ConvertImage converts a Go image.Image to SVG string
func ConvertImage(img image.Image, config *Config) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Convert to RGBA byte array
	pixels := make([]byte, width*height*4)
	idx := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[idx] = byte(r >> 8)
			pixels[idx+1] = byte(g >> 8)
			pixels[idx+2] = byte(b >> 8)
			pixels[idx+3] = byte(a >> 8)
			idx += 4
		}
	}

	return ConvertBytes(pixels, width, height, config)
}

// ConvertBytes converts RGBA byte array to SVG string
func ConvertBytes(pixels []byte, width, height int, config *Config) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if len(pixels) != width*height*4 {
		return "", fmt.Errorf("vtracer: pixel data length mismatch, expected %d, got %d", width*height*4, len(pixels))
	}

	cConfig := config.toCConfig()
	var output *C.char

	result := C.vtracer_convert_bytes(
		(*C.uchar)(unsafe.Pointer(&pixels[0])),
		C.size_t(len(pixels)),
		C.size_t(width),
		C.size_t(height),
		&cConfig,
		&output,
	)

	if err := codeToError(int(result)); err != nil {
		return "", err
	}

	defer C.vtracer_free_string(output)
	return C.GoString(output), nil
}
