# 定义变量
BINARY_NAME=lemon_tea_local
BUILD_DIR=build
LOCAL_PATH=cmd/local/main.go

# 定义所有支持的平台
PLATFORMS=macos_amd64 macos_arm64 windows_amd64 windows_arm64 linux_amd64 linux_arm64

# 创建构建目录
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# 构建相关目标
.PHONY: build
build: _handle_build

.PHONY: _handle_build
_handle_build:
	@if [ "$(filter-out build,$(MAKECMDGOALS))" = "all" ]; then \
		$(MAKE) _build_all; \
	else \
		platform="$(filter-out build,$(MAKECMDGOALS))"; \
		if [ -n "$$platform" ]; then \
			$(MAKE) _build_$$platform; \
		else \
			echo "请指定平台，例如: make build all 或 make build macos_arm64"; \
			exit 1; \
		fi; \
	fi
	@echo "构建完成"

# Proto 相关目标
.PHONY: proto common service all
proto: _handle_proto

.PHONY: _handle_proto
_handle_proto:
	@if [ "$(filter-out proto,$(MAKECMDGOALS))" = "common" ]; then \
		$(MAKE) _proto_common; \
	elif [ "$(filter-out proto,$(MAKECMDGOALS))" = "service" ]; then \
		$(MAKE) _proto_service; \
	elif [ "$(filter-out proto,$(MAKECMDGOALS))" = "all" ]; then \
		$(MAKE) _proto_all; \
	else \
		echo "请指定 proto 文件，例如: make proto common, make proto service 或 make proto all"; \
		exit 1; \
	fi

# 为 proto 参数创建空目标
common service all:
	@:

# 生成 common.proto 的 Go 代码
.PHONY: _proto_common
_proto_common:
	@echo "生成 common.proto 的 Go 代码..."
	protoc --experimental_allow_proto3_optional --go_out=. --openapiv2_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative rpc/common/common.proto
	@echo "生成完成"

# 生成 service 目录下所有 proto 文件的 Go 代码
.PHONY: _proto_service
_proto_service:
	@echo "生成 service 目录下所有 proto 文件的 Go 代码..."
	@if [ ! -d "/tmp/googleapis" ]; then \
		echo "正在下载 googleapis proto 文件..."; \
		git clone https://github.com/googleapis/googleapis.git /tmp/googleapis; \
	fi
	protoc --experimental_allow_proto3_optional -I. -I/tmp/googleapis --openapiv2_out=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative --go-grpc_opt=paths=source_relative rpc/service/*.proto
	@echo "生成完成"

# 生成所有 proto 文件的 Go 代码
.PHONY: _proto_all
_proto_all:
	@echo "生成所有 proto 文件的 Go 代码..."
	@$(MAKE) _proto_common
	@$(MAKE) _proto_service
	@echo "所有 proto 文件生成完成"

# 内部目标 - 构建所有平台
.PHONY: _build_all
_build_all: $(BUILD_DIR)
	@echo "正在构建所有平台..."
	@$(MAKE) _build_macos_amd64
	@$(MAKE) _build_macos_arm64
	@$(MAKE) _build_windows_amd64
	@$(MAKE) _build_windows_arm64
	@$(MAKE) _build_linux_amd64
	@$(MAKE) _build_linux_arm64

# 针对不同平台的编译规则 - 已弃用，请使用 make build {平台}
.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@echo "警告: make $@ 已弃用，请使用 make build $@"
	@exit 1

# 内部编译规则 - 不直接调用
.PHONY: _build_macos_amd64
_build_macos_amd64: $(BUILD_DIR)
	@echo "构建 macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_amd64 $(LOCAL_PATH)

.PHONY: _build_macos_arm64
_build_macos_arm64: $(BUILD_DIR)
	@echo "构建 macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_arm64 $(LOCAL_PATH)

.PHONY: _build_windows_amd64
_build_windows_amd64: $(BUILD_DIR)
	@echo "构建 Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_windows_amd64.exe $(LOCAL_PATH)

.PHONY: _build_windows_arm64
_build_windows_arm64: $(BUILD_DIR)
	@echo "构建 Windows (arm64)..."
	GOOS=windows GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_windows_arm64.exe $(LOCAL_PATH)

.PHONY: _build_linux_amd64
_build_linux_amd64: $(BUILD_DIR)
	@echo "构建 Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 $(LOCAL_PATH)

.PHONY: _build_linux_arm64
_build_linux_arm64: $(BUILD_DIR)
	@echo "构建 Linux (arm64)..."
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64 $(LOCAL_PATH)

# 清理构建产物
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
