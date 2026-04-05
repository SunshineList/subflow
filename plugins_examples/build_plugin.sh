#!/bin/bash

set -e

# ===== 🌊 清理残留文件 =====
rm -f go.mod go.sum Dockerfile

# ===== 🧩 参数检查 =====
if [ -z "$1" ]; then
    echo "❌ 用法: $0 <plugin_go_file>"
    exit 1
fi

PLUGIN_FILE="$1"
PLUGIN_NAME="${PLUGIN_FILE%.go}"  # 去掉 .go 后缀
SO_FILE="${PLUGIN_NAME}.so"
TMP_IMAGE="plugin-builder-${PLUGIN_NAME}"
TMP_CONTAINER="plugin_tmp_${PLUGIN_NAME}"

echo "🧩 插件源码: $PLUGIN_FILE"
echo "📦 插件输出: $SO_FILE"

# ===== 🌐 下载主项目 go.mod 和 go.sum =====
echo "🌐 下载 go.mod 和 go.sum..."
wget -q -O go.mod https://raw.githubusercontent.com/SunshineList/subflow/main/go.mod
wget -q -O go.sum https://raw.githubusercontent.com/SunshineList/subflow/main/go.sum

# ===== 🔍 检查插件文件是否存在 =====
if [ ! -f "$PLUGIN_FILE" ]; then
    echo "❌ 文件 $PLUGIN_FILE 不存在"
    exit 1
fi

# ===== 🛠 生成 Dockerfile =====
echo "🛠 生成 Dockerfile..."

cat > Dockerfile <<EOF
FROM golang:1.24.3 AS backend-builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ${PLUGIN_FILE} .

RUN go mod tidy

RUN go build -buildmode=plugin -o ${SO_FILE} ${PLUGIN_FILE}

FROM alpine:latest AS export-stage
WORKDIR /plugin
COPY --from=backend-builder /app/${SO_FILE} .
EOF

# ===== 🐳 构建 Docker 镜像 =====
echo "🐳 构建 Docker 镜像..."
sudo docker build --no-cache -f Dockerfile -t ${TMP_IMAGE} .

# ===== 📦 创建容器并提取插件 =====
echo "📦 创建临时容器并提取插件..."
sudo docker rm -f ${TMP_CONTAINER} 2>/dev/null || true
sudo docker create --name ${TMP_CONTAINER} ${TMP_IMAGE}
sudo docker cp ${TMP_CONTAINER}:/plugin/${SO_FILE} ./${SO_FILE}
sudo docker rm ${TMP_CONTAINER}

echo "✅ 插件已输出到 ./${SO_FILE}"

# ===== 🧹 清理构建文件 =====
rm -f Dockerfile go.mod go.sum
sudo docker rmi -f ${TMP_IMAGE} > /dev/null 2>&1 || true
