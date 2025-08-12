# VTracer Go 绑定 Makefile

# 检测操作系统
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
    LIB_EXT = dylib
    LIB_PREFIX = lib
endif
ifeq ($(UNAME_S),Linux)
    LIB_EXT = so
    LIB_PREFIX = lib
endif
ifeq ($(findstring MINGW,$(UNAME_S)),MINGW)
    LIB_EXT = dll
    LIB_PREFIX = 
endif

LIB_NAME = $(LIB_PREFIX)vtracer_go.$(LIB_EXT)

.PHONY: all build-rust build-examples test clean install help

# 默认目标
all: test

# 构建示例程序
build-examples:
	@echo "🔨 构建示例程序..."
	@mkdir -p bin
	cd examples/simple && go build -o ../../bin/simple_convert .
	cd examples/advanced && go build -o ../../bin/advanced_convert .
	cd examples/batch && go build -o ../../bin/batch_convert .
	cd examples/web && go build -o ../../bin/http_server .
	@echo "✅ 示例程序构建完成"

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v
	@echo "✅ 测试完成"

# 运行性能测试
bench:
	@echo "📊 运行性能测试..."
	go test -bench=. -benchmem
	@echo "✅ 性能测试完成"

# 清理生成文件
clean:
	@echo "🧹 清理文件..."
	cargo clean
	rm -f $(LIB_NAME)
	rm -rf bin/
	rm -f test_*.png test_*.svg
	@echo "✅ 清理完成"

# 安装到系统（需要 sudo）
install:
	@echo "📦 安装库文件到系统..."
ifeq ($(UNAME_S),Darwin)
	sudo cp $(LIB_NAME) /usr/local/lib/
	sudo install_name_tool -id /usr/local/lib/$(LIB_NAME) /usr/local/lib/$(LIB_NAME)
endif
ifeq ($(UNAME_S),Linux)
	sudo cp $(LIB_NAME) /usr/local/lib/
	sudo ldconfig
endif
	@echo "✅ 安装完成"

# 快速示例测试
demo: build-examples
	@echo "🎨 运行演示..."
	@echo "创建测试图像..."
	@go run examples/create_test_image.go
	@echo "转换测试图像..."
	./bin/simple_convert test_input.png test_output.svg
	@echo "✅ 演示完成，检查 test_output.svg"

# 运行 Web 服务器
serve: build-examples
	@echo "🌐 启动 Web 服务器..."
	./bin/http_server

# 检查依赖
check-deps:
	@echo "🔍 检查依赖..."
	@command -v cargo >/dev/null 2>&1 || { echo "❌ 需要安装 Rust"; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "❌ 需要安装 Go"; exit 1; }
	@echo "✅ 依赖检查通过"

# 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	cargo fmt
	go fmt ./...
	@echo "✅ 代码格式化完成"

# 显示帮助信息
help:
	@echo "VTracer Go 绑定构建工具"
	@echo ""
	@echo "可用命令:"
	@echo "  make              - 构建库并运行测试"
	@echo "  make build-rust   - 只构建 Rust 动态库"
	@echo "  make build-examples - 构建所有示例程序"
	@echo "  make test         - 运行测试"
	@echo "  make bench        - 运行性能测试"
	@echo "  make demo         - 运行快速演示"
	@echo "  make serve        - 启动 Web 服务器"
	@echo "  make clean        - 清理生成文件"
	@echo "  make install      - 安装到系统 (需要 sudo)"
	@echo "  make check-deps   - 检查依赖"
	@echo "  make fmt          - 格式化代码"
	@echo "  make help         - 显示此帮助"
	@echo ""
	@echo "示例用法:"
	@echo "  make && ./bin/simple_convert image.jpg output.svg"
	@echo "  make serve  # 然后访问 http://localhost:8080"
