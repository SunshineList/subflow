# 1. 构建前端
FROM node:20-alpine AS frontend-builder
WORKDIR /frontend
COPY webs ./webs
RUN npm install -g pnpm
RUN cd webs && pnpm install && pnpm run build

# 2. 构建后端（利用 Go 原生交叉编译，不依赖 QEMU）
FROM --platform=$BUILDPLATFORM golang:1.24.3 AS backend-builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /frontend/webs/dist ./static
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -tags=prod -ldflags="-s -w" -o subflow .

# 3. 运行镜像
FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends tzdata ca-certificates wget && \
    ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /app/db /app/logs /app/template /app/plugins && \
    chmod 777 /app/db /app/logs /app/template /app/plugins

COPY --from=backend-builder /app/subflow /app/subflow
COPY --from=backend-builder /app/static /app/static

EXPOSE 8000
CMD ["/app/subflow"]
