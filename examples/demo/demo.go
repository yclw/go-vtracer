package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/yclw/go-vtracer"
)

func main() {
	fmt.Println("🎨 VTracer Go 绑定演示")
	fmt.Println()

	// 1. 创建一个简单的测试图像
	fmt.Println("1. 创建测试图像...")
	img := createTestImage(200, 200)

	// 保存测试图像
	saveImage(img, "demo_input.png")
	fmt.Println("   ✅ 测试图像已保存: demo_input.png")

	// 2. 使用默认配置转换
	fmt.Println("\n2. 使用默认配置转换...")
	err := vtracer.ConvertFile("demo_input.png", "demo_default.svg", nil)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}
	fmt.Println("   ✅ 默认配置转换完成: demo_default.svg")

	// 3. 使用不同预设
	presets := []struct {
		preset vtracer.Preset
		name   string
		file   string
	}{
		{vtracer.PresetBW, "黑白模式", "demo_bw.svg"},
		{vtracer.PresetPoster, "海报模式", "demo_poster.svg"},
		{vtracer.PresetPhoto, "照片模式", "demo_photo.svg"},
	}

	fmt.Println("\n3. 使用不同预设转换...")
	for _, p := range presets {
		config := vtracer.NewConfigFromPreset(p.preset)
		err := vtracer.ConvertFile("demo_input.png", p.file, config)
		if err != nil {
			log.Printf("   ❌ %s 转换失败: %v", p.name, err)
		} else {
			fmt.Printf("   ✅ %s 转换完成: %s\n", p.name, p.file)
		}
	}

	// 4. 直接处理 Go image.Image
	fmt.Println("\n4. 处理 Go image.Image...")
	gradientImg := createGradientImage(150, 150)
	svg, err := vtracer.ConvertImage(gradientImg, nil)
	if err != nil {
		log.Printf("   ❌ 图像转换失败: %v", err)
	} else {
		err = os.WriteFile("demo_gradient.svg", []byte(svg), 0644)
		if err != nil {
			log.Printf("   ❌ 保存失败: %v", err)
		} else {
			fmt.Printf("   ✅ 渐变图像转换完成: demo_gradient.svg (长度: %d 字符)\n", len(svg))
		}
	}

	// 5. 自定义配置
	fmt.Println("\n5. 使用自定义配置...")
	customConfig := &vtracer.Config{
		ColorMode:       vtracer.ColorModeColor,
		Mode:            vtracer.PathModeSpline,
		FilterSpeckle:   8,
		ColorPrecision:  8,
		LayerDifference: 32,
		CornerThreshold: 90,
		LengthThreshold: 3.5,
		MaxIterations:   15, // 必须 > 0
		SpliceThreshold: 30,
		PathPrecision:   3,
	}

	err = vtracer.ConvertFile("demo_input.png", "demo_custom.svg", customConfig)
	if err != nil {
		log.Printf("   ❌ 自定义配置转换失败: %v", err)
	} else {
		fmt.Println("   ✅ 自定义配置转换完成: demo_custom.svg")
	}

	fmt.Println("\n🎉 演示完成！生成的文件:")
	files := []string{
		"demo_input.png",
		"demo_default.svg",
		"demo_bw.svg",
		"demo_poster.svg",
		"demo_photo.svg",
		"demo_gradient.svg",
		"demo_custom.svg",
	}

	for _, file := range files {
		if info, err := os.Stat(file); err == nil {
			fmt.Printf("   📄 %s (%d 字节)\n", file, info.Size())
		}
	}

	fmt.Println()
	fmt.Println("💡 提示:")
	fmt.Println("   - 用浏览器打开 .svg 文件查看结果")
	fmt.Println("   - 比较不同预设的效果")
	fmt.Println("   - 尝试修改自定义配置参数")
}

// 创建测试图像
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 背景白色
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	// 添加一些彩色形状
	centerX, centerY := width/2, height/2

	// 红色圆形
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx, dy := x-centerX/2, y-centerY/2
			if dx*dx+dy*dy < (width/8)*(width/8) {
				img.Set(x, y, color.RGBA{255, 100, 100, 255})
			}
		}
	}

	// 蓝色矩形
	for y := centerY; y < height*3/4; y++ {
		for x := centerX; x < width*3/4; x++ {
			img.Set(x, y, color.RGBA{100, 100, 255, 255})
		}
	}

	// 绿色三角形
	for y := height / 4; y < centerY; y++ {
		for x := centerX + (y - height/4); x < width*3/4-(y-height/4); x++ {
			if x >= 0 && x < width {
				img.Set(x, y, color.RGBA{100, 255, 100, 255})
			}
		}
	}

	return img
}

// 创建渐变图像
func createGradientImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 彩虹渐变
			r := uint8(255 * x / width)
			g := uint8(255 * y / height)
			b := uint8(255 * (x + y) / (width + height))
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

// 保存图像
func saveImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
