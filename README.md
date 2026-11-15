# install-software

一个 Go 编写的命令行安装工具：**安装时自动解析并下载软件的最新版本**，交付只需要一个二进制文件。

## 安装 install-software 自身

从 [Releases](https://github.com/gamelife2025/install-software/releases/latest) 下载二进制：

```bash
wget https://github.com/gamelife2025/install-software/releases/latest/download/install-software-linux-amd64 -O install-software
chmod +x install-software
sudo mv install-software /usr/local/bin/
```

升级自身：

```bash
install-software upgrade
```

## 使用

```bash
install-software <app> [版本]   # 安装软件，版本可选（部分应用支持指定版本）
install-software list           # 列出所有支持的软件
install-software help           # 帮助
```

需要 root 权限的操作（写 /opt、/usr/local/bin、/usr/share/applications）会自动使用 `sudo`。

### desktop

```text
install-software dbeaver        或 install-software 1
install-software drawio         或 install-software 2
install-software firefox        或 install-software 3
install-software sublime_text   或 install-software 4
install-software vscode         或 install-software 5
install-software pycharm        或 install-software 6
install-software joplin         或 install-software 7
install-software qq             或 install-software 8
install-software qqmusic        或 install-software 9
install-software redisInsight   或 install-software 10
install-software hoppscotch     或 install-software 11
```

### cli

```text
install-software docker-compose
install-software etcd
install-software gh
install-software go
install-software hugo
install-software java
install-software protoc
install-software syncthing
install-software upx
install-software vscode-cli
install-software nvm
```

指定版本示例：

```bash
install-software go 1.27.4
install-software hugo v0.170.0
```

> GitHub API 匿名请求有 60 次/小时的限额，可设置 `GITHUB_TOKEN` 环境变量提高限额。

## 从源码构建

```bash
go build -o install-software .
./install-software help
```

## 发布

推送 `v*` tag 会触发 GitHub Actions 自动构建并上传 `install-software-linux-amd64`（含 sha256 校验文件）到 Release：

```bash
git tag v1.0.0 && git push origin v1.0.0
```

## 原理

每个软件对应一个「解析器 + 安装器」：

- **解析器**在安装时实时获取最新版本：GitHub Releases API（gh/hugo/etcd 等）、官网重定向（firefox/vscode）、厂商 API（JetBrains/go.dev）、页面正则（QQ/QQ 音乐）
- **安装器**统一下载到临时目录，解压/拷贝到 `/opt` 或 `/usr/local/bin`，并安装内嵌的 desktop 文件和图标（`go:embed`）

不再需要在仓库里硬编码版本号，也不需要定时任务更新版本。

## 目录结构

```text
main.go / selfupdate.go   入口调度、自升级
core/                     公共层：App 定义、注册表、HTTP/下载/解压/sudo、内嵌资源
  desktop/  icon/         内嵌的 desktop 文件和图标
apps/                     每个软件一个文件，init() 自动注册
  app_<名称>.go           桌面应用（如 app_vscode.go）
  app_cli_<名称>.go       命令行工具（如 app_cli_go.go）
```

新增一个软件只需在 `apps/` 下添加一个文件，无需改动其他代码。
