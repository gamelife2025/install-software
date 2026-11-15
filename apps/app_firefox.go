package apps

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(firefoxApp)
}

var firefoxRe = regexp.MustCompile(`firefox-(.*)\.tar`)

func firefoxURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://download-installer.cdn.mozilla.net/pub/firefox/releases/%s/linux-x86_64/zh-CN/firefox-%s.tar.xz", v, v), true
}

var firefoxApp = core.App{
	Name:     "firefox",
	Alias:    "3",
	Category: "desktop",
	Latest: func() (string, string, error) {
		url, err := core.ResolveRedirect("https://download.mozilla.org/?product=firefox-latest-ssl&os=linux64&lang=zh-CN")
		if err != nil {
			return "", "", err
		}
		m := firefoxRe.FindStringSubmatch(url)
		if m == nil {
			return "", "", fmt.Errorf("无法从下载地址解析 firefox 版本: %s", url)
		}
		return m[1], url, nil
	},
	PinURL: firefoxURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, filepath.Base(url))
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.ExtractTarRoot(file, "/opt", 0); err != nil {
			return err
		}
		if err := core.InstallDesktopFile("firefox.desktop", ""); err != nil {
			return err
		}
		// desktop 文件同时放到 /opt/firefox 下，保持原有布局
		data, _ := core.ReadAsset("desktop/firefox.desktop")
		tmp := filepath.Join(dir, "firefox.desktop")
		if err := os.WriteFile(tmp, data, 0644); err != nil {
			return err
		}
		if err := core.SudoRun("cp", tmp, "/opt/firefox/firefox.desktop"); err != nil {
			return err
		}
		if err := core.SudoRun("ln", "-fs", "/opt/firefox/firefox", "/usr/local/bin/firefox"); err != nil {
			return err
		}
		fmt.Println("install firefox latest success")
		return nil
	},
}
