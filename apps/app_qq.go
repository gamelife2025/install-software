package apps

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(qqApp)
}

var qqRe = regexp.MustCompile(`https://dldir1\.qq\.com/qqfile/qq/QQNT/Linux/QQ_(\d+\.\d+\.\d+)_[0-9]{1,6}_x86_\d+_\d+\.AppImage`)

var qqApp = core.App{
	Name:     "qq",
	Alias:    "8",
	Category: "desktop",
	Latest: func() (string, string, error) {
		resp, err := core.HTTPGet("https://cdn-go.cn/qq-web/im.qq.com_new/latest/rainbow/linuxQQDownload.js")
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}
		m := qqRe.FindSubmatch(body)
		if m == nil {
			return "", "", fmt.Errorf("没有提取到 qq 下载地址")
		}
		return string(m[1]), string(m[0]), nil
	},
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "qq.AppImage")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.InstallAppImage(file, "qq", "qq.png", "qq.desktop", ""); err != nil {
			return err
		}
		fmt.Printf("install qq %s success\n", version)
		return nil
	},
}
