# Makefile for go_ImagePreviewer

APP_NAME = go_ImagePreviewer
BUILD_FLAGS = -ldflags="-s -w" -trimpath

# Windows用（コンソールなし）
windows:
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -ldflags="-s -w -H=windowsgui" -o $(APP_NAME).exe

# Windows用（コンソール付き、デバッグ用）
windows-console:
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(APP_NAME)_console.exe

# Linux用
linux:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(APP_NAME)

# macOS用
macos:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(APP_NAME)

# 現在のプラットフォーム用
build:
	go build $(BUILD_FLAGS)

# 全プラットフォーム用
all: windows linux macos

# クリーンアップ
clean:
	rm -f $(APP_NAME) $(APP_NAME).exe $(APP_NAME)_console.exe

.PHONY: windows windows-console linux macos build all clean
