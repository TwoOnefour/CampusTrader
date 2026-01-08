
.PHONY: build-ui build-go pack upload restart clean

# 变量定义
UI_DIR = frontend
ASSETS_DIR = internal/assets/dist

# 1. 完整构建：先做前端，再做后端
build: build-ui build-go pack build-docker

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
	CGO_ENABLED=1 go build -ldflags="-s -w -extldflags '-static'" -o CampusTrader cmd/main.go

pack:
	rm -rf campustrader.tar.gz
	tar -zcvf campustrader.tar.gz ./CampusTrader ./.env ./static ./sqlite.db

build-docker:
	docker build . -t twoonefour1/campustrader
	docker push twoonefour1/campustrader
