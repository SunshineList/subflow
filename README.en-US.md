<div align="center">

<img src="webs/src/assets/logo.png" width="120px" height="120px" />

# SubFlow

🚀 **All-in-one Proxy Subscription Conversion & Management Platform**

Effortlessly manage all your proxy nodes and convert between client formats

[![Go](https://img.shields.io/badge/Go-1.24.3-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.4-4FC08D?logo=vue.js&logoColor=white)](https://vuejs.org/)
[![Element Plus](https://img.shields.io/badge/Element%20Plus-2.6-409EFF?logo=element&logoColor=white)](https://element-plus.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Docker Image](https://img.shields.io/badge/GHCR-ready-2496ED?logo=docker&logoColor=white)](https://ghcr.io/sunshinelist/subflow)
[![Release](https://img.shields.io/github/v/release/SunshineList/subflow?color=blue)](https://github.com/SunshineList/subflow/releases)

[中文](README.md) | English

</div>

---

## ✨ Why SubFlow?

<table>
<tr>
<td width="50%">

### 🎨 Redesigned Modern UI
- SaaS-style interface with deep Element Plus customization
- Glassmorphism login, gradient cards, smooth animations
- Data visualization dashboard (protocol distribution, access stats, task status)
- Dark mode support, responsive layout

</td>
<td width="50%">

### 🔄 Smart Subscription Management
- Import external subscription links with auto node parsing
- Cron-based scheduled updates to keep nodes fresh
- **Auto-numbering for duplicate node names** — no more update failures
- Drag-and-drop sorting for flexible node ordering

</td>
</tr>
<tr>
<td>

### 📡 Full Protocol Coverage
- **V2Ray** — Base64 universal format
- **Clash** — SS / SSR / Trojan / VMess / VLESS / Hysteria / Hysteria2 / TUIC / AnyTLS / Socks5
- **Surge** — SS / Trojan / VMess / Hysteria2 / TUIC
- Clash `dialer-proxy` chained proxy support

</td>
<td>

### 🧩 Plugin Extension System
- Hot-pluggable architecture, no core code changes needed
- Web UI for plugin management and configuration
- Complete development examples and build scripts
- API event listeners for custom logic

</td>
</tr>
<tr>
<td>

### 🔐 Secure & Self-hosted
- Token-authorized API access
- Independent API Key management
- Full subscription access logging
- Self-hosted, your data stays with you

</td>
<td>

### ⚡ Ready Out of the Box
- Docker one-click deploy, multi-arch (amd64 / arm64)
- One-click install script with auto service registration
- SQLite zero-config database
- Automated CI/CD build & release

</td>
</tr>
</table>

## 📸 Screenshots

| Dashboard | Subscription Management |
|:---:|:---:|
| ![Dashboard](webs/src/assets/1.png) | ![Subscriptions](webs/src/assets/2.png) |

| Node Management | Template Editor |
|:---:|:---:|
| ![Nodes](webs/src/assets/3.png) | ![Templates](webs/src/assets/4.png) |

| Plugin System | Login Page |
|:---:|:---:|
| ![Plugins](webs/src/assets/5.png) | ![Login](webs/src/assets/6.png) |

## 🚀 Quick Start

### Docker Compose (Recommended)

```bash
mkdir subflow && cd subflow
wget https://raw.githubusercontent.com/SunshineList/subflow/main/docker-compose.yml
docker compose up -d
```

Visit `http://your-server-ip:8000` — Default credentials: `admin` / `123456`

### Docker Run

```bash
docker run --name subflow -p 8000:8000 \
  -v $PWD/db:/app/db \
  -v $PWD/template:/app/template \
  -v $PWD/logs:/app/logs \
  -v $PWD/plugins:/app/plugins \
  -d ghcr.io/sunshinelist/subflow
```

### One-click Install (Linux)

```bash
wget https://raw.githubusercontent.com/SunshineList/subflow/main/install.sh && sh install.sh
```

> ⚠️ Alpine Linux uses `musl` instead of `glibc`, so plugins won't work. Docker deployment is recommended.

## 🧩 Plugin Development

Three steps to extend SubFlow:

```bash
# 1. Reference the example plugin
cat plugins_examples/email_plugin.go

# 2. Build the plugin
wget https://raw.githubusercontent.com/SunshineList/subflow/main/plugins_examples/build_plugin.sh
chmod +x build_plugin.sh && ./build_plugin.sh your_plugin.go

# 3. Deploy and enable via Web UI
cp your_plugin.so plugins/
```

<details>
<summary>📋 Plugin Interface</summary>

```go
type Plugin interface {
    Name() string
    Version() string
    Description() string
    DefaultConfig() map[string]interface{}
    SetConfig(map[string]interface{})
    Init() error
    Close() error
    OnAPIEvent(ctx *gin.Context, event EventType, path string,
               statusCode int, requestBody interface{},
               responseBody interface{}) error
    InterestedAPIs() []string
    InterestedEvents() []EventType
}
```

</details>

## 🏗️ Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Vue 3 + Element Plus + Vite + Pinia + TypeScript |
| Backend | Go + Gin + Gorm + SQLite |
| Deployment | Docker (multi-arch) + GitHub Actions CI/CD |
| Package Mgr | pnpm (frontend) + Go Modules (backend) |

## 🙏 Acknowledgements

This project is based on the following open-source projects:

- [eun1e/sublinkE](https://github.com/eun1e/sublinkE) — Original project, thanks to eun1e for the hard work
- [sublinkX](https://github.com/gooaclok819/sublinkX) — Upstream project of sublinkE
- [vue3-element-admin](https://github.com/youlaitech/vue3-element-admin) — Frontend template framework

## 📄 License

[MIT License](LICENSE)
