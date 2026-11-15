package apps

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(qqmusicApp)
}

// 页面中的下载链接是带签名的重定向地址，直接构造的 dldir1 地址已 404
var qqmusicRe = regexp.MustCompile(`https://c\.y\.qq\.com/cgi-bin/file_redirect\.fcg[^"']*qqmusic-(\d+\.\d+\.\d+)\.AppImage[^"']*`)

var qqmusicApp = core.App{
	Name:     "qqmusic",
	Alias:    "9",
	Category: "desktop",
	Latest: func() (string, string, error) {
		resp, err := core.HTTPGet("https://y.qq.com/download/download.html")
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}
		m := qqmusicRe.FindSubmatch(body)
		if m == nil {
			return "", "", fmt.Errorf("没有提取到 qqmusic 下载地址")
		}
		return string(m[1]), string(m[0]), nil
	},
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "qqmusic.AppImage")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.InstallAppImage(file, "qqmusic", "qqmusic.png", "qqmusic.desktop", ""); err != nil {
			return err
		}
		fmt.Printf("install qqmusic %s success\n", version)
		return nil
	},
}
