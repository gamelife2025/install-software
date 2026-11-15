package apps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(upxApp)
}

func upxURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://github.com/upx/upx/releases/download/v%s/upx-%s-amd64_linux.tar.xz", v, v), true
}

var upxApp = core.App{
	Name:     "upx",
	Category: "cli",
	Latest: func() (string, string, error) {
		tag, err := core.GithubLatestTag("upx/upx")
		if err != nil {
			return "", "", err
		}
		url, _ := upxURL(tag)
		return tag, url, nil
	},
	PinURL: upxURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "upx.tar.xz")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.ExtractTarUser(file, dir); err != nil {
			return err
		}
		bin := filepath.Join(dir, fmt.Sprintf("upx-%s-amd64_linux/upx", core.TrimV(version)))
		if err := core.SudoRun("cp", bin, "/usr/local/bin/upx"); err != nil {
			return err
		}
		fmt.Printf("install upx %s success\n", version)
		return nil
	},
}
