package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yclw/go-vtracer"
)

func main() {
	var (
		input      = flag.String("input", "", "输入图像文件路径")
		output     = flag.String("output", "", "输出 SVG 文件路径")
		preset     = flag.String("preset", "", "预设配置 (bw|poster|photo)")
		colorMode  = flag.String("colormode", "color", "颜色模式 (color|binary)")
		mode       = flag.String("mode", "spline", "路径模式 (pixel|polygon|spline)")
		speckle    = flag.Int("speckle", 4, "过滤斑点大小")
		precision  = flag.Int("precision", 6, "颜色精度 (1-8)")
		difference = flag.Int("difference", 16, "图层差异 (0-255)")
		corner     = flag.Int("corner", 60, "角度阈值 (0-180)")
		length     = flag.Float64("length", 4.0, "长度阈值 (3.5-10.0)")
		splice     = flag.Int("splice", 45, "拼接阈值 (0-180)")
		pathPrec   = flag.Int("pathprec", 2, "路径精度")
	)

	flag.Parse()

	if *input == "" || *output == "" {
		fmt.Println("VTracer Go 绑定 - 高级转换工具")
		fmt.Println()
		flag.Usage()
		fmt.Println()
		fmt.Println("预设选项:")
		fmt.Println("  bw     - 黑白模式，适用于简单图形")
		fmt.Println("  poster - 海报模式，适用于颜色较少的插图")
		fmt.Println("  photo  - 照片模式，适用于复杂的照片")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  advanced_convert -input photo.jpg -output photo.svg -preset photo")
		fmt.Println("  advanced_convert -input logo.png -output logo.svg -preset poster")
		fmt.Println("  advanced_convert -input sketch.jpg -output sketch.svg -colormode binary")
		os.Exit(1)
	}

	// 检查输入文件
	if _, err := os.Stat(*input); os.IsNotExist(err) {
		log.Fatalf("输入文件不存在: %s", *input)
	}

	var config *vtracer.Config

	// 使用预设配置
	if *preset != "" {
		switch strings.ToLower(*preset) {
		case "bw":
			config = vtracer.NewConfigFromPreset(vtracer.PresetBW)
			fmt.Println("使用黑白预设配置")
		case "poster":
			config = vtracer.NewConfigFromPreset(vtracer.PresetPoster)
			fmt.Println("使用海报预设配置")
		case "photo":
			config = vtracer.NewConfigFromPreset(vtracer.PresetPhoto)
			fmt.Println("使用照片预设配置")
		default:
			log.Fatalf("未知的预设: %s", *preset)
		}
	} else {
		// 使用自定义配置
		config = vtracer.DefaultConfig()
		fmt.Println("使用自定义配置")

		// 设置颜色模式
		switch strings.ToLower(*colorMode) {
		case "binary", "bw":
			config.ColorMode = vtracer.ColorModeBinary
		case "color":
			config.ColorMode = vtracer.ColorModeColor
		default:
			log.Fatalf("未知的颜色模式: %s", *colorMode)
		}

		// 设置路径模式
		switch strings.ToLower(*mode) {
		case "pixel", "none":
			config.Mode = vtracer.PathModeNone
		case "polygon":
			config.Mode = vtracer.PathModePolygon
		case "spline":
			config.Mode = vtracer.PathModeSpline
		default:
			log.Fatalf("未知的路径模式: %s", *mode)
		}

		// 设置其他参数
		config.FilterSpeckle = *speckle
		config.ColorPrecision = *precision
		config.LayerDifference = *difference
		config.CornerThreshold = *corner
		config.LengthThreshold = *length
		config.SpliceThreshold = *splice
		config.PathPrecision = *pathPrec

		// 验证参数范围
		if config.ColorPrecision < 1 || config.ColorPrecision > 8 {
			log.Fatal("颜色精度必须在 1-8 之间")
		}
		if config.LayerDifference < 0 || config.LayerDifference > 255 {
			log.Fatal("图层差异必须在 0-255 之间")
		}
		if config.CornerThreshold < 0 || config.CornerThreshold > 180 {
			log.Fatal("角度阈值必须在 0-180 之间")
		}
		if config.LengthThreshold < 3.5 || config.LengthThreshold > 10.0 {
			log.Fatal("长度阈值必须在 3.5-10.0 之间")
		}
		if config.SpliceThreshold < 0 || config.SpliceThreshold > 180 {
			log.Fatal("拼接阈值必须在 0-180 之间")
		}
	}

	// 显示配置信息
	fmt.Printf("配置信息:\n")
	fmt.Printf("  颜色模式: %s\n", map[vtracer.ColorMode]string{
		vtracer.ColorModeColor:  "彩色",
		vtracer.ColorModeBinary: "二值",
	}[config.ColorMode])
	fmt.Printf("  路径模式: %s\n", map[vtracer.PathMode]string{
		vtracer.PathModeNone:    "像素",
		vtracer.PathModePolygon: "多边形",
		vtracer.PathModeSpline:  "样条",
	}[config.Mode])
	fmt.Printf("  过滤斑点: %d\n", config.FilterSpeckle)
	fmt.Printf("  颜色精度: %d\n", config.ColorPrecision)
	fmt.Printf("  图层差异: %d\n", config.LayerDifference)

	fmt.Printf("正在转换 %s -> %s\n", *input, *output)

	// 执行转换
	err := vtracer.ConvertFile(*input, *output, config)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}

	// 显示结果信息
	if info, err := os.Stat(*output); err == nil {
		fmt.Printf("转换成功！输出文件大小: %d 字节\n", info.Size())
	} else {
		fmt.Println("转换成功！")
	}
}
