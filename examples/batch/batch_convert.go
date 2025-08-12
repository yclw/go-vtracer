package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/yclw/go-vtracer"
)

func main() {
	var (
		inputDir   = flag.String("input", "", "输入目录")
		outputDir  = flag.String("output", "", "输出目录")
		preset     = flag.String("preset", "poster", "预设配置 (bw|poster|photo)")
		workers    = flag.Int("workers", 4, "并发处理数量")
		extensions = flag.String("ext", "jpg,jpeg,png,bmp,tiff", "支持的文件扩展名（逗号分隔）")
		recursive  = flag.Bool("recursive", false, "递归处理子目录")
	)

	flag.Parse()

	if *inputDir == "" || *outputDir == "" {
		fmt.Println("VTracer Go 绑定 - 批量转换工具")
		fmt.Println()
		flag.Usage()
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  batch_convert -input ./images -output ./svg -preset photo")
		fmt.Println("  batch_convert -input ./photos -output ./vectors -preset photo -workers 8")
		os.Exit(1)
	}

	// 检查输入目录
	if _, err := os.Stat(*inputDir); os.IsNotExist(err) {
		log.Fatalf("输入目录不存在: %s", *inputDir)
	}

	// 创建输出目录
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	// 解析文件扩展名
	extList := strings.Split(*extensions, ",")
	for i, ext := range extList {
		extList[i] = "." + strings.ToLower(strings.TrimSpace(ext))
	}

	// 获取配置
	var config *vtracer.Config
	switch strings.ToLower(*preset) {
	case "bw":
		config = vtracer.NewConfigFromPreset(vtracer.PresetBW)
	case "poster":
		config = vtracer.NewConfigFromPreset(vtracer.PresetPoster)
	case "photo":
		config = vtracer.NewConfigFromPreset(vtracer.PresetPhoto)
	default:
		log.Fatalf("未知的预设: %s", *preset)
	}

	fmt.Printf("批量转换配置:\n")
	fmt.Printf("  输入目录: %s\n", *inputDir)
	fmt.Printf("  输出目录: %s\n", *outputDir)
	fmt.Printf("  预设配置: %s\n", *preset)
	fmt.Printf("  并发数量: %d\n", *workers)
	fmt.Printf("  支持扩展名: %v\n", extList)
	fmt.Printf("  递归处理: %t\n", *recursive)
	fmt.Println()

	// 收集要处理的文件
	var files []string
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if !*recursive && path != *inputDir {
				return filepath.SkipDir
			}
			return nil
		}

		// 检查文件扩展名
		ext := strings.ToLower(filepath.Ext(path))
		for _, supportedExt := range extList {
			if ext == supportedExt {
				files = append(files, path)
				break
			}
		}

		return nil
	}

	if err := filepath.Walk(*inputDir, walkFunc); err != nil {
		log.Fatalf("扫描文件失败: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("未找到支持的图像文件")
		return
	}

	fmt.Printf("找到 %d 个文件需要处理\n", len(files))

	// 创建工作队列
	fileQueue := make(chan string, len(files))
	for _, file := range files {
		fileQueue <- file
	}
	close(fileQueue)

	// 结果统计
	var (
		successCount int
		failureCount int
		totalSize    int64
		mu           sync.Mutex
		wg           sync.WaitGroup
	)

	// 启动工作协程
	startTime := time.Now()
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for inputPath := range fileQueue {
				// 计算输出路径
				relPath, err := filepath.Rel(*inputDir, inputPath)
				if err != nil {
					log.Printf("Worker %d: 计算相对路径失败 %s: %v", workerID, inputPath, err)
					continue
				}

				outputPath := filepath.Join(*outputDir, relPath)
				outputPath = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".svg"

				// 确保输出目录存在
				outputDirPath := filepath.Dir(outputPath)
				if err := os.MkdirAll(outputDirPath, 0755); err != nil {
					log.Printf("Worker %d: 创建输出目录失败 %s: %v", workerID, outputDirPath, err)
					mu.Lock()
					failureCount++
					mu.Unlock()
					continue
				}

				// 转换文件
				fmt.Printf("Worker %d: 处理 %s\n", workerID, filepath.Base(inputPath))
				err = vtracer.ConvertFile(inputPath, outputPath, config)

				mu.Lock()
				if err != nil {
					log.Printf("Worker %d: 转换失败 %s: %v", workerID, inputPath, err)
					failureCount++
				} else {
					successCount++
					if info, statErr := os.Stat(outputPath); statErr == nil {
						totalSize += info.Size()
					}
				}
				mu.Unlock()
			}
		}(i + 1)
	}

	// 等待所有工作完成
	wg.Wait()
	duration := time.Since(startTime)

	// 输出统计信息
	fmt.Println()
	fmt.Println("=== 批量转换完成 ===")
	fmt.Printf("处理时间: %v\n", duration.Round(time.Millisecond))
	fmt.Printf("成功转换: %d 个文件\n", successCount)
	fmt.Printf("转换失败: %d 个文件\n", failureCount)
	fmt.Printf("总输出大小: %.2f MB\n", float64(totalSize)/(1024*1024))

	if successCount > 0 {
		avgTime := duration / time.Duration(successCount)
		fmt.Printf("平均转换时间: %v/文件\n", avgTime.Round(time.Millisecond))
	}
}
