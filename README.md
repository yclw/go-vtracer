# VTracer Go 绑定

[![Go Reference](https://pkg.go.dev/badge/github.com/yclw/go-vtracer.svg)](https://pkg.go.dev/github.com/yclw/go-vtracer)
[![Go Report Card](https://goreportcard.com/badge/github.com/yclw/go-vtracer)](https://goreportcard.com/report/github.com/yclw/go-vtracer)
[![License](https://img.shields.io/badge/license-MIT%2FApache--2.0-blue.svg)](LICENSE)

这是 [VTracer](https://github.com/visioncortex/vtracer) 的 **Go 语言绑定**，让您能在 Go 程序中轻松使用 VTracer 进行高质量的图像矢量化处理。

## 功能特性

- **高质量矢量化**：将光栅图像转换为紧凑的 SVG 矢量图
- **高性能**：直接调用优化的 Rust 核心库，比纯 Go 实现快 10-100 倍
- **易于使用**：符合 Go 语言习惯的简洁 API
- **多种输入**：支持文件路径、Go image.Image、字节数组
- **灵活配置**：完全自定义参数
- **跨平台**：支持 Linux (x86_64/aarch64)、macOS (x86_64/arm64)、Windows (x86_64)

## 快速开始

### 1. 安装包

```bash
go get github.com/yclw/go-vtracer
```

### 2. 开始使用

```go
package main

import (
    "log"
    "github.com/yclw/go-vtracer"
)

func main() {
    // 使用默认配置转换
    err := vtracer.ConvertFile("input.jpg", "output.svg", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // 或者使用自定义配置
    config := vtracer.DefaultConfig()
    config.ColorPrecision = 8
    config.FilterSpeckle = 8
    err = vtracer.ConvertFile("photo.jpg", "photo.svg", config)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. 运行示例

```bash
go run examples/simple/simple_convert.go input.jpg output.svg
```

## 📁 示例程序

本项目提供了一个简单的示例程序，位于 `examples/simple/` 目录下，展示了最核心的文件转换功能：

```bash
go run examples/simple/simple_convert.go input.jpg output.svg
```

## 使用示例

### 基础文件转换

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yclw/go-vtracer"
)

func main() {
    // 使用默认配置
    err := vtracer.ConvertFile("input.jpg", "output.svg", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("转换成功！")
}
```

### 自定义配置

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yclw/go-vtracer"
)

func main() {
    config := vtracer.DefaultConfig()
    config.ColorMode = vtracer.ColorModeColor
    config.Mode = vtracer.PathModeSpline
    config.FilterSpeckle = 8
    config.ColorPrecision = 6
    config.LayerDifference = 32
    
    err := vtracer.ConvertFile("image.png", "output.svg", config)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("自定义配置转换成功！")
}
```

### 处理 Go image.Image

```go
package main

import (
    "fmt"
    "image"
    "image/png"
    "log"
    "os"
    
    "github.com/yclw/go-vtracer"
)

func main() {
    // 从文件加载图像
    file, err := os.Open("input.png")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    
    img, err := png.Decode(file)
    if err != nil {
        log.Fatal(err)
    }
    
    // 转换为 SVG
    svg, err := vtracer.ConvertImage(img, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // 保存 SVG
    err = os.WriteFile("output.svg", []byte(svg), 0644)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("生成的 SVG 长度: %d 字符\n", len(svg))
}
```

## API 参考

### 类型定义

```go
type ColorMode uint8
const (
    ColorModeColor  ColorMode = 0  // 彩色模式
    ColorModeBinary ColorMode = 1  // 二值模式
)

type Hierarchical uint8
const (
    HierarchicalStacked Hierarchical = 0  // 堆叠模式
    HierarchicalCutout  Hierarchical = 1  // 裁剪模式
)

type PathMode uint8
const (
    PathModeNone    PathMode = 0  // 像素模式
    PathModePolygon PathMode = 1  // 多边形模式
    PathModeSpline  PathMode = 2  // 样条模式
)


```

### 配置选项

```go
type Config struct {
    ColorMode        ColorMode   // 颜色模式
    Hierarchical     Hierarchical // 层次结构模式
    FilterSpeckle    int         // 过滤斑点大小 (像素)
    ColorPrecision   int         // 颜色精度 (1-8 位)
    LayerDifference  int         // 图层颜色差异 (0-255)
    Mode             PathMode    // 路径拟合模式
    CornerThreshold  int         // 角度阈值 (0-180 度)
    LengthThreshold  float64     // 长度阈值 (3.5-10.0)
    MaxIterations    int         // 最大迭代次数
    SpliceThreshold  int         // 拼接阈值 (0-180 度)
    PathPrecision    int         // 路径精度 (小数位数)
}
```

### 主要函数

```go
// 创建默认配置
func DefaultConfig() *Config

// 文件转换
func ConvertFile(inputPath, outputPath string, config *Config) error

// Go 图像转换
func ConvertImage(img image.Image, config *Config) (string, error)

// 字节数组转换
func ConvertBytes(pixels []byte, width, height int, config *Config) (string, error)
```

## 许可证

MIT
