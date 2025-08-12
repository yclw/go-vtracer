package vtracer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Fatal("默认配置不应为 nil")
	}

	fmt.Printf("默认配置: %+v\n", config)
}

// TestPresetConfigs 测试预设配置
func TestPresetConfigs(t *testing.T) {
	presets := []Preset{
		PresetBW,
		PresetPoster,
		PresetPhoto,
	}

	for _, preset := range presets {
		config := NewConfigFromPreset(preset)
		if config == nil {
			t.Fatalf("预设 %d 的配置不应为 nil", preset)
		}
		fmt.Printf("预设 %d 配置: %+v\n", preset, config)
	}
}

// TestCreateSolidColorImage 测试创建纯色图像
func TestCreateSolidColorImage(t *testing.T) {
	img := CreateSolidColorImage(100, 100, color.RGBA{255, 0, 0, 255})
	if img == nil {
		t.Fatal("创建的图像不应为 nil")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Fatalf("期望图像大小为 100x100，实际为 %dx%d", bounds.Dx(), bounds.Dy())
	}
}

// TestConvertSolidColorImage 测试转换纯色图像
func TestConvertSolidColorImage(t *testing.T) {
	// 创建一个简单的红色图像
	img := CreateSolidColorImage(50, 50, color.RGBA{255, 0, 0, 255})

	// 使用默认配置转换
	svg, err := ConvertImage(img, nil)
	if err != nil {
		t.Fatalf("转换图像失败: %v", err)
	}

	if len(svg) == 0 {
		t.Fatal("生成的 SVG 为空")
	}

	fmt.Printf("生成的 SVG 长度: %d 字符\n", len(svg))
	fmt.Printf("SVG 前100字符: %.100s...\n", svg)
}

// TestConvertImageWithDifferentPresets 测试不同预设的转换效果
func TestConvertImageWithDifferentPresets(t *testing.T) {
	// 创建一个更复杂的图像（渐变色）
	img := createGradientImage(100, 100)

	presets := []Preset{
		PresetBW,
		PresetPoster,
		PresetPhoto,
	}

	presetNames := []string{"黑白", "海报", "照片"}

	for i, preset := range presets {
		config := NewConfigFromPreset(preset)
		svg, err := ConvertImage(img, config)

		if err != nil {
			t.Fatalf("使用 %s 预设转换失败: %v", presetNames[i], err)
		}

		if len(svg) == 0 {
			t.Fatalf("使用 %s 预设生成的 SVG 为空", presetNames[i])
		}

		fmt.Printf("%s 预设生成的 SVG 长度: %d 字符\n", presetNames[i], len(svg))
	}
}

// TestConvertFile 测试文件转换（需要实际图像文件）
func TestConvertFile(t *testing.T) {
	// 创建测试图像文件
	testImagePath := "test_image.png"
	testSVGPath := "test_output.svg"

	// 创建测试图像
	img := createTestPattern(200, 200)
	file, err := os.Create(testImagePath)
	if err != nil {
		t.Fatalf("创建测试图像文件失败: %v", err)
	}
	defer os.Remove(testImagePath)
	defer os.Remove(testSVGPath)

	err = png.Encode(file, img)
	file.Close()
	if err != nil {
		t.Fatalf("编码 PNG 失败: %v", err)
	}

	// 转换文件
	config := NewConfigFromPreset(PresetPoster)
	err = ConvertFile(testImagePath, testSVGPath, config)
	if err != nil {
		t.Fatalf("文件转换失败: %v", err)
	}

	// 检查输出文件是否存在
	if _, err := os.Stat(testSVGPath); os.IsNotExist(err) {
		t.Fatal("输出 SVG 文件不存在")
	}

	// 读取并检查 SVG 内容
	svgContent, err := os.ReadFile(testSVGPath)
	if err != nil {
		t.Fatalf("读取 SVG 文件失败: %v", err)
	}

	if len(svgContent) == 0 {
		t.Fatal("SVG 文件为空")
	}

	fmt.Printf("文件转换成功，SVG 大小: %d 字节\n", len(svgContent))
}

// createGradientImage 创建渐变图像
func createGradientImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 创建从红到蓝的渐变
			r := uint8(255 * x / width)
			g := uint8(255 * y / height)
			b := uint8(255 - r)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

// createTestPattern 创建测试图案
func createTestPattern(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 创建一些几何形状
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var c color.Color

			// 背景白色
			c = color.RGBA{255, 255, 255, 255}

			// 添加一些彩色矩形
			if x > width/4 && x < 3*width/4 && y > height/4 && y < 3*height/4 {
				if x < width/2 {
					c = color.RGBA{255, 0, 0, 255} // 红色
				} else {
					c = color.RGBA{0, 0, 255, 255} // 蓝色
				}
			}

			// 添加圆形
			centerX, centerY := width/2, height/2
			dx, dy := x-centerX, y-centerY
			if dx*dx+dy*dy < (width/8)*(width/8) {
				c = color.RGBA{0, 255, 0, 255} // 绿色圆形
			}

			img.Set(x, y, c)
		}
	}

	return img
}

// BenchmarkConvertImage 性能基准测试
func BenchmarkConvertImage(b *testing.B) {
	img := createTestPattern(100, 100)
	config := DefaultConfig()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ConvertImage(img, config)
		if err != nil {
			b.Fatalf("转换失败: %v", err)
		}
	}
}
