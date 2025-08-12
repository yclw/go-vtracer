#!/bin/bash

# VTracer Go 绑定多平台构建脚本

set -e

VERSION=${1:-"v0.1.0"}
RELEASE_DIR="releases"

echo "🚀 开始构建 VTracer Go 绑定 $VERSION"
echo "📂 构建目录: $RELEASE_DIR"

# 清理并创建目录
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

# 定义平台配置
declare -a TARGETS=(
    "x86_64-unknown-linux-gnu:linux-x86_64:libvtracer_go.so"
    "aarch64-unknown-linux-gnu:linux-aarch64:libvtracer_go.so"
    "x86_64-apple-darwin:macos-x86_64:libvtracer_go.dylib"
    "aarch64-apple-darwin:macos-aarch64:libvtracer_go.dylib"
    "x86_64-pc-windows-gnu:windows-x86_64:vtracer_go.dll"
)

# 检查工具
echo "🔧 检查构建工具..."
if ! command -v cargo &> /dev/null; then
    echo "❌ 错误: 未安装 Rust。请访问 https://rustup.rs/ 安装"
    exit 1
fi

if ! command -v zip &> /dev/null; then
    echo "❌ 错误: 未安装 zip"
    exit 1
fi

# 安装目标平台
echo "📦 安装目标平台..."
for target_info in "${TARGETS[@]}"; do
    IFS=':' read -r target name lib_name <<< "$target_info"
    echo "  - 安装 $target"
    rustup target add "$target" || echo "    已安装 $target"
done

# 针对 Linux aarch64 安装交叉编译工具
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "🔧 检查 Linux 交叉编译工具..."
    if ! dpkg -l | grep -q gcc-aarch64-linux-gnu; then
        echo "📥 安装 aarch64 交叉编译工具..."
        sudo apt-get update
        sudo apt-get install -y gcc-aarch64-linux-gnu
    fi
fi

# 构建每个平台
for target_info in "${TARGETS[@]}"; do
    IFS=':' read -r target name lib_name <<< "$target_info"
    
    echo ""
    echo "🔨 构建 $name ($target)..."
    
    # 设置环境变量
    export_cmd=""
    if [[ "$target" == "aarch64-unknown-linux-gnu" ]]; then
        export CARGO_TARGET_AARCH64_UNKNOWN_LINUX_GNU_LINKER=aarch64-linux-gnu-gcc
        export CC_aarch64_unknown_linux_gnu=aarch64-linux-gnu-gcc
        export AR_aarch64_unknown_linux_gnu=aarch64-linux-gnu-ar
    fi
    
    # 构建
    if cargo build --release --target "$target"; then
        echo "  ✅ 构建成功"
        
        # 创建发布包
        lib_path="target/$target/release/$lib_name"
        if [[ -f "$lib_path" ]]; then
            mkdir -p "$RELEASE_DIR/$name"
            cp "$lib_path" "$RELEASE_DIR/$name/"
            
            # 创建 ZIP 包
            cd "$RELEASE_DIR"
            zip -r "vtracer-$target.zip" "$name/"
            cd ..
            
            echo "  📦 创建了 $RELEASE_DIR/vtracer-$target.zip"
        else
            echo "  ❌ 构建文件不存在: $lib_path"
        fi
    else
        echo "  ❌ 构建失败: $target"
        
        # 对于无法交叉编译的平台，提供说明
        if [[ "$target" == *"windows"* ]] && [[ "$OSTYPE" != "msys" ]]; then
            echo "    💡 Windows 构建需要在 Windows 系统或使用 mingw-w64"
        elif [[ "$target" == *"darwin"* ]] && [[ "$OSTYPE" != "darwin"* ]]; then
            echo "    💡 macOS 构建需要在 macOS 系统"
        fi
    fi
done

echo ""
echo "📋 构建完成！生成的文件："
ls -la "$RELEASE_DIR"/*.zip 2>/dev/null || echo "  没有成功构建的包"

echo ""
echo "🚀 下一步："
echo "  1. 检查 $RELEASE_DIR/ 目录中的 ZIP 文件"
echo "  2. 在 GitHub 上创建 release: https://github.com/yclw/go-vtracer/releases/new"
echo "  3. 上传 ZIP 文件作为 release assets"
echo "  4. 设置 tag 为: $VERSION"

echo ""
echo "💡 或者使用 GitHub CLI 自动创建 release："
echo "  gh release create $VERSION $RELEASE_DIR/*.zip --title \"VTracer Go Binding $VERSION\" --notes \"预编译动态库包\""
