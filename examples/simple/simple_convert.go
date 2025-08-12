package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yclw/go-vtracer"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("使用方法: %s <输入图片> <输出SVG>\n", os.Args[0])
		fmt.Println("示例: simple_convert image.jpg output.svg")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// 检查输入文件是否存在
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Fatalf("输入文件不存在: %s", inputPath)
	}

	fmt.Printf("正在转换 %s -> %s\n", inputPath, outputPath)

	// 使用默认配置转换
	err := vtracer.ConvertFile(inputPath, outputPath, nil)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}

	// 检查输出文件大小
	if info, err := os.Stat(outputPath); err == nil {
		fmt.Printf("转换成功！输出文件大小: %d 字节\n", info.Size())
	} else {
		fmt.Println("转换成功！")
	}
}
