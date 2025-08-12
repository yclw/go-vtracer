//go:build ignore

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	releaseURL = "https://github.com/yclw/go-vtracer/releases/latest/download"
)

func main() {
	fmt.Println("🔧 正在安装 VTracer 动态库...")

	// 检测平台
	platform := detectPlatform()
	if platform == "" {
		fmt.Printf("❌ 不支持的平台: %s_%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(1)
	}

	// 创建 lib 目录
	libDir := "lib"
	if err := os.MkdirAll(libDir, 0755); err != nil {
		fmt.Printf("❌ 创建目录失败: %v\n", err)
		os.Exit(1)
	}

	// 下载动态库
	libFile := getLibFileName(platform)
	libPath := filepath.Join(libDir, libFile)

	// 检查是否已存在
	if _, err := os.Stat(libPath); err == nil {
		fmt.Printf("✅ 动态库已存在: %s\n", libPath)
		return
	}

	// 优先尝试本地构建
	fmt.Printf("🔨 平台: %s\n", platform)

	// 尝试本地构建
	if err := buildLocalLibrary(); err != nil {
		fmt.Printf("⚠️  本地构建失败: %v\n", err)

		// 回退到下载预编译包
		downloadURL := fmt.Sprintf("%s/vtracer-%s.zip", releaseURL, platform)
		fmt.Printf("📥 尝试下载预编译包: %s\n", downloadURL)

		if err := downloadAndExtract(downloadURL, libDir, platform); err != nil {
			fmt.Printf("❌ 下载也失败: %v\n", err)
			fmt.Println("\n💡 解决方案:")
			fmt.Println("1. 确保安装了 Rust: https://rustup.rs/")
			fmt.Println("2. 手动构建:")
			fmt.Println("   chmod +x build.sh")
			fmt.Println("   ./build.sh")
			fmt.Println("3. 或者等待项目发布预编译包")
			os.Exit(1)
		}
	}

	fmt.Printf("✅ 安装完成: %s\n", libPath)
}

func buildLocalLibrary() error {
	// 检查是否有 Rust 工具链
	if _, err := exec.LookPath("cargo"); err != nil {
		return fmt.Errorf("未找到 Rust 工具链，请安装 Rust: https://rustup.rs/")
	}

	// 检查是否有 Cargo.toml
	if _, err := os.Stat("Cargo.toml"); err != nil {
		return fmt.Errorf("未找到 Cargo.toml，请确保在项目根目录运行")
	}

	fmt.Println("🦀 开始构建 Rust 动态库...")

	// 执行 cargo build --release
	cmd := exec.Command("cargo", "build", "--release")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cargo build 失败: %v", err)
	}

	// 复制动态库文件
	return copyBuiltLibrary()
}

func copyBuiltLibrary() error {
	var srcPath, destPath string

	switch runtime.GOOS {
	case "darwin":
		srcPath = "target/release/libvtracer_go.dylib"
		destPath = "lib/libvtracer_go.dylib"
	case "linux":
		srcPath = "target/release/libvtracer_go.so"
		destPath = "lib/libvtracer_go.so"
	case "windows":
		srcPath = "target/release/vtracer_go.dll"
		destPath = "lib/vtracer_go.dll"
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	// 检查源文件是否存在
	if _, err := os.Stat(srcPath); err != nil {
		return fmt.Errorf("构建的动态库不存在: %s", srcPath)
	}

	// 创建目标目录
	if err := os.MkdirAll("lib", 0755); err != nil {
		return fmt.Errorf("创建 lib 目录失败: %v", err)
	}

	// 复制文件
	return copyFile(srcPath, destPath)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func detectPlatform() string {
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return "x86_64-unknown-linux-gnu"
		case "arm64":
			return "aarch64-unknown-linux-gnu"
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return "x86_64-apple-darwin"
		case "arm64":
			return "aarch64-apple-darwin"
		}
	case "windows":
		if runtime.GOARCH == "amd64" {
			return "x86_64-pc-windows-gnu"
		}
	}
	return ""
}

func getLibFileName(platform string) string {
	switch {
	case runtime.GOOS == "windows":
		return "vtracer_go.dll"
	case runtime.GOOS == "darwin":
		return "libvtracer_go.dylib"
	case runtime.GOOS == "linux":
		return "libvtracer_go.so"
	}
	return ""
}

func downloadAndExtract(url, destDir, platform string) error {
	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	// 保存到临时文件
	tmpFile, err := os.CreateTemp("", "vtracer-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return err
	}

	// 解压 zip
	if err := tmpFile.Close(); err != nil {
		return err
	}

	return extractZip(tmpFile.Name(), destDir, platform)
}

func extractZip(zipPath, destDir, platform string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	libFile := getLibFileName(platform)
	for _, f := range r.File {
		if f.Name == libFile || filepath.Base(f.Name) == libFile {
			return extractFile(f, filepath.Join(destDir, libFile))
		}
	}

	return fmt.Errorf("在 zip 中未找到库文件: %s", libFile)
}

func extractFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	return err
}
