# VTracer Go 绑定

[![Go Reference](https://pkg.go.dev/badge/github.com/yclw/go-vtracer.svg)](https://pkg.go.dev/github.com/yclw/go-vtracer)
[![Go Report Card](https://goreportcard.com/badge/github.com/yclw/go-vtracer)](https://goreportcard.com/report/github.com/yclw/go-vtracer)
[![License](https://img.shields.io/badge/license-MIT%2FApache--2.0-blue.svg)](LICENSE)

这是 [VTracer](https://github.com/visioncortex/vtracer) 的 **Go 语言绑定**，让您能在 Go 程序中轻松使用 VTracer 进行高质量的图像矢量化处理。

## ✨ 功能特性

- 🎨 **高质量矢量化**：将光栅图像转换为紧凑的 SVG 矢量图
- 🚀 **高性能**：直接调用优化的 Rust 核心库，比纯 Go 实现快 10-100 倍
- 🔧 **易于使用**：符合 Go 语言习惯的简洁 API
- 📦 **多种输入**：支持文件路径、Go image.Image、字节数组
- ⚙️ **灵活配置**：预设配置 + 完全自定义参数
- 🌍 **跨平台**：支持 Linux (x86_64/aarch64)、macOS (x86_64/arm64)、Windows (x86_64)
- 🔒 **内存安全**：自动管理 C 内存，防止内存泄漏

## 🚀 快速开始

### 1. 安装包

```bash
go get github.com/yclw/go-vtracer
```

### 2. 安装动态库

```bash
# 在你的项目目录中执行
go generate github.com/yclw/go-vtracer
```

> **说明**：这个包需要 Rust 编译的动态库。`go generate` 会自动下载适合你系统的预编译库。

### 3. 开始使用

```go
package main

import (
    "log"
    "github.com/yclw/go-vtracer"
)

func main() {
    // 简单转换
    err := vtracer.ConvertFile("input.jpg", "output.svg", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用照片预设获得更好效果
    config := vtracer.NewConfigFromPreset(vtracer.PresetPhoto)
    err = vtracer.ConvertFile("photo.jpg", "photo.svg", config)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 4. 运行演示

```bash
go run examples/demo/demo.go  # 查看各种配置的效果
```

## 📁 示例程序

本项目提供了多个示例程序，位于 `examples/` 目录下：

- **`simple/`** - 简单文件转换示例
- **`advanced/`** - 高级配置和命令行工具
- **`batch/`** - 批量处理示例
- **`web/`** - Web 服务器示例
- **`demo/`** - 完整功能演示

每个示例都可以独立运行：

```bash
go run examples/simple/simple_convert.go input.jpg output.svg
go run examples/web/http_server.go  # 启动 Web 服务
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

### 使用预设配置

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yclw/go-vtracer"
)

func main() {
    // 使用照片预设
    config := vtracer.NewConfigFromPreset(vtracer.PresetPhoto)
    
    err := vtracer.ConvertFile("photo.jpg", "photo.svg", config)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("照片转换成功！")
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
    "image/color"
    "log"
    "os"
    
    "github.com/yclw/go-vtracer"
)

func main() {
    // 创建测试图像
    img := vtracer.CreateSolidColorImage(200, 200, color.RGBA{255, 100, 50, 255})
    
    // 转换为 SVG
    svg, err := vtracer.ConvertImage(img, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // 保存 SVG
    err = os.WriteFile("generated.svg", []byte(svg), 0644)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("生成的 SVG 长度: %d 字符\n", len(svg))
}
```

### Web 服务示例

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "strconv"
    
    "github.com/yclw/go-vtracer"
)

func convertHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
        return
    }
    
    // 读取上传的图像
    file, _, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "读取图像失败", http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    imageData, err := io.ReadAll(file)
    if err != nil {
        http.Error(w, "读取图像数据失败", http.StatusBadRequest)
        return
    }
    
    // 解析配置参数
    preset := r.FormValue("preset")
    var config *vtracer.Config
    
    switch preset {
    case "bw":
        config = vtracer.NewConfigFromPreset(vtracer.PresetBW)
    case "poster":
        config = vtracer.NewConfigFromPreset(vtracer.PresetPoster)
    case "photo":
        config = vtracer.NewConfigFromPreset(vtracer.PresetPhoto)
    default:
        config = vtracer.DefaultConfig()
    }
    
    // 这里需要先将 imageData 解码为 image.Image，然后转换
    // 为简化示例，假设我们有解码函数
    img, err := decodeImage(imageData)
    if err != nil {
        http.Error(w, "解码图像失败", http.StatusBadRequest)
        return
    }
    
    // 转换为 SVG
    svg, err := vtracer.ConvertImage(img, config)
    if err != nil {
        http.Error(w, "转换失败: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 返回 SVG
    w.Header().Set("Content-Type", "image/svg+xml")
    w.Header().Set("Content-Length", strconv.Itoa(len(svg)))
    w.Write([]byte(svg))
}

func main() {
    http.HandleFunc("/convert", convertHandler)
    fmt.Println("服务器启动在 :8080")
    http.ListenAndServe(":8080", nil)
}

// 实际使用中需要实现这个函数
func decodeImage(data []byte) (image.Image, error) {
    // 使用 image.Decode 或相关库解码图像
    return nil, nil
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

type Preset int
const (
    PresetBW     Preset = iota  // 黑白预设
    PresetPoster                // 海报预设
    PresetPhoto                 // 照片预设
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

// 从预设创建配置
func NewConfigFromPreset(preset Preset) *Config

// 文件转换
func ConvertFile(inputPath, outputPath string, config *Config) error

// Go 图像转换
func ConvertImage(img image.Image, config *Config) (string, error)

// 字节数组转换
func ConvertBytes(pixels []byte, width, height int, config *Config) (string, error)

// 创建纯色图像（用于测试）
func CreateSolidColorImage(width, height int, c color.Color) image.Image
```

## 📚 API 参考

### 主要函数

```go
// 文件转换
func ConvertFile(inputPath, outputPath string, config *Config) error

// 内存图像转换
func ConvertImage(img image.Image, config *Config) (string, error)

// 字节数据转换
func ConvertBytes(pixels []byte, width, height int, config *Config) (string, error)

// 配置创建
func DefaultConfig() *Config
func NewConfigFromPreset(preset Preset) *Config
```

### 预设配置

- **PresetBW**: 黑白模式，适合简单图形
- **PresetPoster**: 海报模式，适合插图
- **PresetPhoto**: 照片模式，适合复杂图像

## ❓ 常见问题

**Q: 为什么需要 `go generate`？**
A: 这个包依赖 Rust 编译的动态库，`go generate` 会自动下载适合你系统的预编译库。

**Q: 支持哪些图像格式？**
A: 支持 PNG、JPEG、GIF、BMP、TIFF 等常见格式。

**Q: 如何处理大图像？**
A: 建议调整 `FilterSpeckle` 参数过滤小细节，或先缩放图像再处理。

**Q: 转换很慢怎么办？**
A: 使用合适的预设配置，避免过高的 `ColorPrecision` 设置。

## 许可证

本项目遵循与 VTracer 相同的许可证条款。
