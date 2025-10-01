# Makefile for fs CLI tool
# 支持多平台构建：linux, darwin (macOS), windows

# =============================================
# 配置区
# =============================================

# 项目基本信息
BINARY_NAME := fs
ORG_PATH    := fs # 替换为你的模块路径
MAIN_PATH   := . # 主包路径

# 版本信息（可通过 git 获取）
VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')

# 输出目录
DIST_DIR := dist

# 构建目标（格式：GOOS-GOARCH）
PLATFORMS := \
	linux-amd64 \
	linux-arm64 \
	darwin-amd64 \
	darwin-arm64 \
	windows-amd64 \
	windows-arm64

# =============================================
# 默认目标
# =============================================

.PHONY: all build clean dist clean-dist

all: build

# 构建当前系统的版本
build:
	@echo "🏗️  构建当前系统版本..."
	@go build -ldflags="-X '${ORG_PATH}/internal/version.Version=${VERSION}' -X '${ORG_PATH}/internal/version.Commit=${COMMIT}' -X '${ORG_PATH}/internal/version.Date=${DATE}'" -o ${BINARY_NAME} ${MAIN_PATH}
	@echo "✅ 构建完成: ./${BINARY_NAME}"

# =============================================
# 跨平台构建
# =============================================

# 一键打包所有平台
dist: clean-dist $(addprefix dist/,$(addsuffix /${BINARY_NAME},$(PLATFORMS)))
	@echo "🎉 所有平台构建完成！产物在 ./${DIST_DIR}/"

# 清理分发目录
clean-dist:
	rm -rf ${DIST_DIR}
	@echo "🧹 已清理 ${DIST_DIR}/"

# 定义每个平台的构建规则
$(addprefix dist/,$(addsuffix /${BINARY_NAME},$(PLATFORMS))): dist/%/${BINARY_NAME}
	@mkdir -p $(dir $@)
	@GOOS=$(word 1, $(subst -, ,$*)) \
	GOARCH=$(word 2, $(subst -, ,$*)) \
	CGO_ENABLED=0 \
	go build \
	-ldflags="-s -w -X '${ORG_PATH}/internal/version.Version=${VERSION}' -X '${ORG_PATH}/internal/version.Commit=${COMMIT}' -X '${ORG_PATH}/internal/version.Date=${DATE}'" \
	-o $@ ${MAIN_PATH}
	@echo "✅ 构建完成: $@"

# =============================================
# 快捷构建单个平台（可选）
# =============================================

.PHONY: linux mac windows

linux:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
	-ldflags="-s -w -X '${ORG_PATH}/internal/version.Version=${VERSION}' -X '${ORG_PATH}/internal/version.Commit=${COMMIT}' -X '${ORG_PATH}/internal/version.Date=${DATE}'" \
	-o ${BINARY_NAME}-linux-amd64 ${MAIN_PATH}
	@echo "✅ 构建完成: ./${BINARY_NAME}-linux-amd64"

mac:
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build \
	-ldflags="-s -w -X '${ORG_PATH}/internal/version.Version=${VERSION}' -X '${ORG_PATH}/internal/version.Commit=${COMMIT}' -X '${ORG_PATH}/internal/version.Date=${DATE}'" \
	-o ${BINARY_NAME}-darwin-amd64 ${MAIN_PATH}
	@echo "✅ 构建完成: ./${BINARY_NAME}-darwin-amd64"

mac-arm64:
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build \
	-ldflags="-s -w -X '${ORG_PATH}/internal/version.Version=${VERSION}' -X '${ORG_PATH}/internal/version.Commit=${COMMIT}' -X '${ORG_PATH}/internal/version.Date=${DATE}'" \
	-o ${BINARY_NAME}-darwin-arm64 ${MAIN_PATH}
	@echo "✅ 构建完成: ./${BINARY_NAME}-darwin-arm64"

windows:
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
	-ldflags="-s -w -X '${ORG_PATH}/internal/version.Version=${VERSION}' -X '${ORG_PATH}/internal/version.Commit=${COMMIT}' -X '${ORG_PATH}/internal/version.Date=${DATE}'" \
	-o ${BINARY_NAME}-windows-amd64.exe ${MAIN_PATH}
	@echo "✅ 构建完成: ./${BINARY_NAME}-windows-amd64.exe"

# =============================================
# 清理
# =============================================

clean:
	rm -f ${BINARY_NAME} ${BINARY_NAME}-*
	@echo "🧹 已清理构建产物"

# =============================================
# 工具：查看版本
# =============================================

.PHONY: version

version:
	@echo "Version: ${VERSION}"
	@echo "Commit:  ${COMMIT}"
	@echo "Date:    ${DATE}"