#!/bin/bash

# spacebattle 构建脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印函数
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Go 是否安装
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go 未安装，请先安装 Go"
        exit 1
    fi
    print_info "Go 版本: $(go version)"
}

# 安装依赖
install_deps() {
    print_info "安装依赖..."
    go mod download
    go mod tidy
}

# 运行测试
run_tests() {
    print_info "运行测试..."
    go test ./...
}

# 构建游戏
build_game() {
    local platform=$1
    local arch=$2
    local output_name=$3
    
    print_info "构建 $platform-$arch 版本..."
    
    if [ "$platform" = "windows" ]; then
        GOOS=$platform GOARCH=$arch go build -o "build/${output_name}.exe" ./cmd/game
    else
        GOOS=$platform GOARCH=$arch go build -o "build/${output_name}" ./cmd/game
    fi
}

# 主函数
main() {
    print_info "开始构建 spacebattle..."
    
    # 检查 Go
    check_go
    
    # 创建构建目录
    mkdir -p build
    
    # 安装依赖
    install_deps
    
    # 运行测试
    run_tests
    
    # 构建不同平台版本
    build_game "windows" "amd64" "spacebattle-windows-amd64"
    build_game "linux" "amd64" "spacebattle-linux-amd64"
    build_game "darwin" "amd64" "spacebattle-darwin-amd64"
    build_game "darwin" "arm64" "spacebattle-darwin-arm64"
    
    print_info "构建完成！"
    print_info "构建文件位于 build/ 目录"
    
    # 显示构建结果
    echo ""
    print_info "构建结果:"
    ls -la build/
}

# 运行主函数
main "$@"
