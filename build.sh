#!/bin/bash

# VTracer Go 绑定构建脚本

set -e

echo "🦀 构建 VTracer Go 绑定..."

# 检查是否安装了 Rust
if ! command -v cargo &> /dev/null; then
    echo "❌ 错误: 未找到 Rust 工具链，请先安装 Rust"
    echo "访问 https://rustup.rs/ 获取安装说明"
    exit 1
fi

# 检查是否安装了 Go
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到 Go 工具链，请先安装 Go"
    echo "访问 https://golang.org/dl/ 获取安装说明"
    exit 1
fi

echo "✅ 工具链检查通过"

# 构建 Rust 动态库
echo "🔨 构建 Rust 动态库..."
cargo build --release

if [ $? -ne 0 ]; then
    echo "❌ Rust 库构建失败"
    exit 1
fi

echo "✅ Rust 库构建成功"

# 复制动态库到当前目录
echo "📁 复制动态库文件..."

# 检测操作系统并复制相应的库文件
case "$(uname -s)" in
    Darwin*)
        if [ -f "target/release/libvtracer_go.dylib" ]; then
            cp target/release/libvtracer_go.dylib .
            echo "✅ 已复制 libvtracer_go.dylib"
        else
            echo "❌ 未找到 libvtracer_go.dylib"
            exit 1
        fi
        ;;
    Linux*)
        if [ -f "target/release/libvtracer_go.so" ]; then
            cp target/release/libvtracer_go.so .
            echo "✅ 已复制 libvtracer_go.so"
        else
            echo "❌ 未找到 libvtracer_go.so"
            exit 1
        fi
        ;;
    MINGW*|MSYS*|CYGWIN*)
        if [ -f "target/release/vtracer_go.dll" ]; then
            cp target/release/vtracer_go.dll .
            echo "✅ 已复制 vtracer_go.dll"
        else
            echo "❌ 未找到 vtracer_go.dll"
            exit 1
        fi
        ;;
    *)
        echo "❌ 不支持的操作系统: $(uname -s)"
        exit 1
        ;;
esac

# 初始化 Go 模块（如果尚未初始化）
if [ ! -f "go.mod" ]; then
    echo "📦 初始化 Go 模块..."
    go mod init github.com/yclw/go-vtracer
    echo "✅ Go 模块初始化完成"
fi

# 运行 Go 测试
echo "🧪 运行 Go 测试..."
go test -v

if [ $? -eq 0 ]; then
    echo "✅ 所有测试通过！"
else
    echo "❌ 测试失败"
    exit 1
fi

echo ""
echo "🎉 VTracer Go 绑定构建完成！"
echo ""
echo "使用方法:"
echo "  import \"github.com/yclw/go-vtracer\""
echo ""
echo "示例:"
echo "  vtracer.ConvertFile(\"input.jpg\", \"output.svg\", nil)"
echo ""
