
.PHONY: build-ui build-go build clean

# 变量定义
UI_DIR = frontend
ASSETS_DIR = internal/assets/dist
BINARY_NAME = CampusTrader

# 1. 完整构建：先做前端，再做后端
build: build-ui build-go

# 2. 编译前端
build-ui:
	@echo "正在构建前端..."
	cd $(UI_DIR) && bash -i -c "nvm use 22 && yarn install && npm run build"
	@echo "拷贝前端产物到 Go 目录..."
	rm -rf $(ASSETS_DIR)/*
	cp -r $(UI_DIR)/dist/* $(ASSETS_DIR)/

# 3. 静态编译 Go 后端
build-go:
	@echo "正在构建 Go 二进制文件 (静态编译)..."
	go build -ldflags="-s -w" -o CampusTrader cmd/main.go

# 4. 清理产物
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(ASSETS_DIR)/*