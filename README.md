<div align="center">

# SubFlow

**开源订阅转换管理工具**

基于 [eun1e/sublinkE](https://github.com/eun1e/sublinkE) 二次开发

[![Go](https://img.shields.io/badge/Go-1.24.3-00ADD8?logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.4-4FC08D?logo=vue.js)](https://vuejs.org/)
[![Element Plus](https://img.shields.io/badge/Element%20Plus-2.6-409EFF)](https://element-plus.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-available-2496ED?logo=docker)](https://hub.docker.com/)

中文 | [English](README.en-US.md)

</div>

---

## 致谢

本项目基于以下优秀开源项目进行二次开发，特别感谢原作者的付出与贡献：

- **[sublinkE](https://github.com/eun1e/sublinkE)** - 原始项目，由 [eun1e](https://github.com/eun1e) 开发维护，感谢 eun1e 的辛勤付出
- **[sublinkX](https://github.com/gooaclok819/sublinkX)** - sublinkE 的上游项目
- **[vue3-element-admin](https://github.com/youlaitech/vue3-element-admin)** - 前端模板框架

## 项目简介

SubFlow 是一个高自由度的开源订阅转换管理工具，支持多种代理协议和客户端格式。

- 前端采用 Vue3 + Element Plus + Vite
- 后端采用 Go + Gin + Gorm + SQLite
- 默认账号：`admin` 密码：`123456`，请安装后务必自行修改

## 功能特性

- 高自由度与安全性，支持访问订阅记录及简易配置管理
- 支持多种客户端协议及格式：
  - **V2Ray**（Base64 通用格式）
  - **Clash**（支持 SS/SSR/Trojan/VMess/VLESS/Hysteria/Hysteria2/TUIC/AnyTLS/Socks5）
  - **Surge**（支持 SS/Trojan/VMess/Hysteria2/TUIC）
- Token 授权访问 API
- 导入、定时更新订阅链接中的节点
- 订阅节点拖拽排序
- Clash `dialer-proxy` 前置代理支持
- API Key 管理
- 插件扩展系统（实验性）
- **节点重名自动编号**（避免订阅更新时重名冲突）
- **现代化 UI 界面**

## 安装部署

### Docker Compose 部署（推荐）

```bash
# 创建并进入项目目录
mkdir subflow && cd subflow

# 下载 docker-compose.yml
wget https://raw.githubusercontent.com/SunshineList/sublinkE/main/docker-compose.yml

# 启动服务
docker compose up -d
```

访问 `http://your-server-ip:8000` 即可使用。

### Docker 运行

```bash
docker run --name subflow -p 8000:8000 \
  -v $PWD/db:/app/db \
  -v $PWD/template:/app/template \
  -v $PWD/logs:/app/logs \
  -v $PWD/plugins:/app/plugins \
  -d sunshinelist/subflow
```

### 一键安装

```bash
wget https://raw.githubusercontent.com/SunshineList/sublinkE/refs/heads/main/install.sh && sh install.sh
```

> **注意**：在 Alpine Linux 上运行一键安装脚本时，由于 Alpine 使用 `musl` 而非 `glibc`，插件模块无法正常工作。推荐优先使用 Docker 部署以获得最佳兼容性。

## 插件系统

SubFlow 提供了灵活的插件系统，允许开发者扩展功能而无需修改核心代码。

### 开发步骤

1. 参照 `plugins_examples/email_plugin.go` 编写自定义插件
2. 使用 `plugins_examples/build_plugin.sh` 编译成 `.so` 文件
3. 将编译好的 `.so` 文件放入 `plugins` 目录

### 插件接口

所有插件必须实现 `plugins.Plugin` 接口：

```go
Name() string
Version() string
Description() string
DefaultConfig() map[string]interface{}
SetConfig(map[string]interface{})
Init() error
Close() error
OnAPIEvent(ctx *gin.Context, event plugins.EventType, path string,
           statusCode int, requestBody interface{},
           responseBody interface{}) error
InterestedAPIs() []string
InterestedEvents() []plugins.EventType
```

### 编译插件

```bash
wget https://raw.githubusercontent.com/SunshineList/sublinkE/main/plugins_examples/build_plugin.sh
chmod +x build_plugin.sh
./build_plugin.sh your_plugin.go
cp your_plugin.so plugins/
```

可通过 Web 界面管理插件的启用/禁用、配置参数和运行状态。

## 项目预览

![预览1](webs/src/assets/1.png)
![预览2](webs/src/assets/2.png)
![预览3](webs/src/assets/3.png)
![预览4](webs/src/assets/4.png)
![预览5](webs/src/assets/5.png)
![预览6](webs/src/assets/6.png)

## License

[MIT License](LICENSE)
