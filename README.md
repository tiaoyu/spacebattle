# spacebattle

基于 [Ebiten](https://github.com/hajimehoshi/ebiten) 游戏引擎开发的 2D 游戏项目。

## 项目结构

```
spacebattle/
├── cmd/game/              # 游戏入口点
│   └── main.go
├── internal/              # 内部包
│   ├── game/             # 游戏核心逻辑
│   ├── scenes/           # 游戏场景
│   ├── entities/         # 游戏实体
│   ├── utils/            # 工具函数
│   └── config/           # 配置管理
├── assets/               # 游戏资源
│   ├── images/           # 图像资源
│   ├── sounds/           # 音频资源
│   └── fonts/            # 字体资源
├── docs/                 # 文档
│   ├── design/           # 设计文档
│   └── api/              # API 文档
├── scripts/              # 构建脚本
├── tests/                # 测试文件
├── go.mod               # Go 模块文件
├── Makefile             # 构建配置
└── README.md            # 项目说明
```

## 快速开始

### 环境要求

- Go 1.21 或更高版本
- 支持 OpenGL 的图形环境

### 安装依赖

```bash
go mod download
```

### 运行游戏

```bash
# 直接运行
go run ./cmd/game

# 或使用 Makefile
make run
```

### 构建游戏

```bash
# 构建当前平台
make build

# 跨平台构建
make build-all

# 使用构建脚本
./scripts/build.sh
```

## 开发指南

### 添加新场景

1. 在 `internal/scenes/` 目录下创建新的场景文件
2. 实现 `Scene` 接口的 `Update()` 和 `Draw()` 方法
3. 在游戏逻辑中注册新场景

### 添加新实体

1. 在 `internal/entities/` 目录下创建新的实体文件
2. 定义实体结构体和相关方法
3. 在场景中使用新实体

### 资源管理

- 图像资源放在 `assets/images/` 目录
- 音频资源放在 `assets/sounds/` 目录
- 字体资源放在 `assets/fonts/` 目录

## 开发命令

```bash
# 运行游戏
make run

# 构建游戏
make build

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make vet

# 清理构建文件
make clean

# 查看所有命令
make help
```

## 设计文档

- [游戏设计文档](docs/design/game_design_document.md)
- [技术设计文档](docs/design/technical_design.md)

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 致谢

- [Ebiten](https://github.com/hajimehoshi/ebiten) - 优秀的 Go 语言 2D 游戏引擎
- [Go](https://golang.org/) - 强大的编程语言
