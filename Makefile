
.PHONY: build-ui build-go pack upload restart clean

# 变量定义
UI_DIR = frontend
ASSETS_DIR = internal/assets/dist
BINARY_NAME = CampusTrader
REMOTE_USER = root
REMOTE_HOST = lucianawa.cn
REMOTE_ADDR = $(REMOTE_USER)@$(REMOTE_HOST)
REMOTE_PORT = 23333
REMOTE_PATH = /opt/campustrader
TARGET_FILE = campustrader.tar.gz

# 1. 完整构建：先做前端，再做后端
build: build-ui build-go pack upload restart clean

# 2. 编译前端
build-ui:
	@echo "正在构建前端..."
	cd $(UI_DIR) && bash -i -c "nvm use 24 && yarn install && npm run build"
	@echo "拷贝前端产物到 Go 目录..."
	rm -rf $(ASSETS_DIR)/*
	mkdir -p $(ASSETS_DIR)
	cp -r $(UI_DIR)/dist/* $(ASSETS_DIR)/

# 3. 静态编译 Go 后端
build-go:
	@echo "正在构建 Go 二进制文件 (静态编译)..."
	go build -ldflags="-s -w" -o CampusTrader cmd/main.go

pack:
	rm -rf campustrader.tar.gz
	tar -zcvf campustrader.tar.gz ./CampusTrader ./.env ./static ./sqlite.db

restart:
	ssh -p $(REMOTE_PORT) $(REMOTE_ADDR) 'cd $(REMOTE_PATH) && \
		tar -zxvf $(TARGET_FILE) && \
		systemctl restart campustrader'

upload:
	scp -P $(REMOTE_PORT) $(TARGET_FILE) $(REMOTE_ADDR):$(REMOTE_PATH)/

# 4. 清理产物
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(ASSETS_DIR)/*
	rm -rf campustrader.tar.gz
	ssh -p $(REMOTE_PORT) $(REMOTE_ADDR) 'cd $(REMOTE_PATH) && rm -rf $(TARGET_FILE)'