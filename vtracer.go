//go:generate go run cmd/install/install.go

package vtracer

/*
#cgo linux LDFLAGS: -L./lib -lvtracer_go -Wl,-rpath,./lib
#cgo darwin LDFLAGS: -L./lib -lvtracer_go -Wl,-rpath,./lib
#cgo windows LDFLAGS: -L./lib -lvtracer_go
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
	"errors"
	"image"
	"image/color"
	"unsafe"
)

// ColorMode 表示颜色模式
type ColorMode uint8

const (
	ColorModeColor  ColorMode = 0
	ColorModeBinary ColorMode = 1
)

// Hierarchical 表示层次结构模式
type Hierarchical uint8

const (
	HierarchicalStacked Hierarchical = 0
	HierarchicalCutout  Hierarchical = 1
)

// PathMode 表示路径模式
type PathMode uint8

const (
	PathModeNone    PathMode = 0
	PathModePolygon PathMode = 1
	PathModeSpline  PathMode = 2
)

// Config 包含 vtracer 的配置选项
type Config struct {
	ColorMode       ColorMode    // 颜色模式：彩色或二值
	Hierarchical    Hierarchical // 层次结构：堆叠或裁剪
	FilterSpeckle   int          // 过滤斑点大小
	ColorPrecision  int          // 颜色精度 (1-8)
	LayerDifference int          // 图层差异 (0-255)
	Mode            PathMode     // 路径拟合模式
	CornerThreshold int          // 角度阈值 (0-180度)
	LengthThreshold float64      // 长度阈值 (3.5-10.0)
	MaxIterations   int          // 最大迭代次数
	SpliceThreshold int          // 拼接阈值 (0-180度)
	PathPrecision   int          // 路径精度
}

// DefaultConfig 返回默认配置
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

// toCConfig 将 Go 配置转换为 C 配置
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

// ConvertFile 将图像文件转换为 SVG 文件
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

	switch result {
	case 0:
		return nil
	case -1:
		return errors.New("无效的输入参数")
	case -2:
		return errors.New("无效的输入路径")
	case -3:
		return errors.New("无效的输出路径")
	case -4:
		return errors.New("转换失败")
	default:
		return errors.New("未知错误")
	}
}

// ConvertImage 将 Go image.Image 转换为 SVG 字符串
func ConvertImage(img image.Image, config *Config) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 转换为 RGBA 字节数组
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

// ConvertBytes 将 RGBA 字节数组转换为 SVG 字符串
func ConvertBytes(pixels []byte, width, height int, config *Config) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if len(pixels) != width*height*4 {
		return "", errors.New("像素数据长度不匹配")
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

	if result != 0 {
		switch result {
		case -1:
			return "", errors.New("无效的输入参数")
		case -2:
			return "", errors.New("像素数据长度错误")
		case -3:
			return "", errors.New("字符串转换失败")
		case -4:
			return "", errors.New("图像转换失败")
		default:
			return "", errors.New("未知错误")
		}
	}

	defer C.vtracer_free_string(output)
	return C.GoString(output), nil
}

// Preset 预设配置
type Preset int

const (
	PresetBW     Preset = iota // 黑白模式
	PresetPoster               // 海报模式
	PresetPhoto                // 照片模式
)

// NewConfigFromPreset 根据预设创建配置
func NewConfigFromPreset(preset Preset) *Config {
	config := DefaultConfig()

	switch preset {
	case PresetBW:
		config.ColorMode = ColorModeBinary
		config.FilterSpeckle = 4
		config.ColorPrecision = 6
		config.LayerDifference = 16
	case PresetPoster:
		config.ColorMode = ColorModeColor
		config.FilterSpeckle = 4
		config.ColorPrecision = 8
		config.LayerDifference = 16
	case PresetPhoto:
		config.ColorMode = ColorModeColor
		config.FilterSpeckle = 10
		config.ColorPrecision = 8
		config.LayerDifference = 48
		config.CornerThreshold = 180
	}

	return config
}

// CreateSolidColorImage 创建纯色图像用于测试
func CreateSolidColorImage(width, height int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}
