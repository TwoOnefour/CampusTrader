# ==========================================
# 阶段 1: 前端构建 (UI Builder)
# 对应 Makefile 中的 build-ui
# ==========================================
FROM node:24-alpine AS ui-builder

# 设置工作目
WORKDIR /app/frontend

# 优化缓存：先拷贝依赖文件
COPY frontend/package.json ./

# 安装依赖
RUN yarn install

# 拷贝前端源代码
COPY frontend .

# 编译前端 (npm run build)
RUN npm run build

# ==========================================
# 阶段 2: 后端构建 (Go Builder)
# 对应 Makefile 中的 build-go
# ==========================================
FROM golang:1.24-alpine AS go-builder

# 安装构建依赖 (CGO_ENABLED=1 需要 gcc 和 musl-dev)
RUN apk add --no-cache build-base

WORKDIR /src

# 优化缓存：先拷贝 Go 依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 拷贝所有后端源码
COPY . .

# 【关键步骤】将阶段1构建好的前端静态资源拷贝到 Go 代码指定的目录
# 对应 Makefile: cp -r $(UI_DIR)/dist/* $(ASSETS_DIR)/
RUN mkdir -p internal/assets/dist && \
    rm -rf internal/assets/dist/*
COPY --from=ui-builder /app/frontend/dist ./internal/assets/dist

# 静态编译 Go 二进制文件
# 对应 Makefile: CGO_ENABLED=1 go build -ldflags="-s -w -extldflags '-static'" ...
RUN CGO_ENABLED=1 go build -ldflags="-s -w -extldflags '-static'" -o CampusTrader cmd/main.go

# ==========================================
# 阶段 3: 最终运行镜像 (Final Runtime)
# ==========================================
FROM alpine:latest

# 安装基础证书 (如果应用需要访问 HTTPS 外部链接) 和 时区数据
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 1. 从构建阶段拷贝编译好的二进制文件
COPY --from=go-builder /src/CampusTrader .

# 2. 拷贝本地的 static 文件夹 (根据你的要求)
COPY static ./static

# 3. 拷贝 sqlite.db (根据你的要求)
# 注意：如果构建时本地目录下没有这个文件，这一步会报错
COPY sqlite.db ./sqlite.db

# 暴露端口 (假设是 8080，请根据实际情况修改)
EXPOSE 8080

# 启动命令
CMD ["./CampusTrader"]
