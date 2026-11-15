package apps

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(vscodeApp)
}

// resolveVSCode 解析 vscode 最新稳定版，返回 commitID 和下载地址
// vscode 和 vscode-cli 共用
func resolveVSCode() (commitID, url string, err error) {
	url, err = core.ResolveRedirect("https://code.visualstudio.com/sha/download?build=stable&os=linux-x64")
	if err != nil {
		return "", "", err
	}
	if !strings.HasSuffix(url, ".tar.gz") {
		return "", "", fmt.Errorf("意料之外的 vscode 下载地址: %s", url)
	}
	// 国内 CDN 镜像加速（旧域名返回时替换）
	url = strings.Replace(url, "az764295.vo.msecnd.net/stable", "vscode.cdn.azure.cn/stable", 1)
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("无法从下载地址解析 vscode commit: %s", url)
	}
	return parts[len(parts)-2], url, nil
}

var vscodeApp = core.App{
	Name:     "vscode",
	Alias:    "5",
	Category: "desktop",
	Latest: func() (string, string, error) {
		commitID, url, err := resolveVSCode()
		return commitID, url, err
	},
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, filepath.Base(url))
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.ExtractTarRoot(file, "/opt", 0); err != nil {
			return err
		}
		if err := core.InstallIcon("vscode.png", "/opt/VSCode-linux-x64"); err != nil {
			return err
		}
		if err := core.InstallDesktopFile("code.desktop", ""); err != nil {
			return err
		}
		if err := core.InstallDesktopFile("code-url-handler.desktop", ""); err != nil {
			return err
		}
		if err := core.SudoRun("ln", "-fs", "/opt/VSCode-linux-x64/code", "/usr/local/bin/code"); err != nil {
			return err
		}
		if err := core.SudoRun("ln", "-fs", "/opt/VSCode-linux-x64/code-tunnel", "/usr/local/bin/code-tunnel"); err != nil {
			return err
		}
		fmt.Println("install vscode latest success")
		return nil
	},
}
