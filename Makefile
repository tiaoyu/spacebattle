# spacebattle Makefile

# 变量定义
BINARY_NAME=spacebattle
BUILD_DIR=build
MAIN_PATH=./cmd/game

# 默认目标
.PHONY: all
all: build

# 构建游戏
.PHONY: build
build:
	@echo "构建游戏..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行游戏
.PHONY: run
run:
	@echo "运行游戏..."
	@go run $(MAIN_PATH)

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@go clean

# 安装依赖
.PHONY: deps
deps:
	@echo "安装依赖..."
	@go mod download
	@go mod tidy

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	@go test ./...

# 运行测试并显示覆盖率
.PHONY: test-coverage
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	@go fmt ./...

# 代码检查
.PHONY: vet
vet:
	@echo "运行代码检查..."
	@go vet ./...

# 跨平台构建
.PHONY: build-all
build-all:
	@echo "跨平台构建..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "跨平台构建完成"

# 开发模式运行（带热重载）
.PHONY: dev
dev:
	@echo "开发模式运行..."
	@air

# 帮助信息
.PHONY: help
help:
	@echo "可用的命令:"
	@echo "  build        - 构建游戏"
	@echo "  run          - 运行游戏"
	@echo "  clean        - 清理构建文件"
	@echo "  deps         - 安装依赖"
	@echo "  test         - 运行测试"
	@echo "  test-coverage- 运行测试并生成覆盖率报告"
	@echo "  fmt          - 格式化代码"
	@echo "  vet          - 代码检查"
	@echo "  build-all    - 跨平台构建"
	@echo "  dev          - 开发模式运行"
	@echo "  help         - 显示帮助信息"
