package apps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(protocApp)
}

func protocURL(version string) (string, bool) {
	v := core.TrimV(version)
	// 3.0.0 以下需要源码编译，不再支持
	if v == "" || !core.VersionNewer("2.9.9", v) {
		return "", false
	}
	return fmt.Sprintf("https://github.com/protocolbuffers/protobuf/releases/download/v%s/protoc-%s-linux-x86_64.zip", v, v), true
}

var protocApp = core.App{
	Name:     "protoc",
	Category: "cli",
	Latest: func() (string, string, error) {
		tag, err := core.GithubLatestTag("protocolbuffers/protobuf")
		if err != nil {
			return "", "", err
		}
		url, ok := protocURL(tag)
		if !ok {
			return "", "", fmt.Errorf("protoc %s 版本过低，仅支持 3.0.0 及以上", tag)
		}
		return tag, url, nil
	},
	PinURL: protocURL,
	Install: func(version, url string) error {
		if _, err := exec.LookPath("unzip"); err != nil {
			return fmt.Errorf("请先安装 unzip")
		}
		file, dir, err := core.FetchToTemp(url, "protoc.zip")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.Run(exec.Command("unzip", "-o", file, "-d", filepath.Join(dir, "protoc"))); err != nil {
			return err
		}
		if err := core.SudoRun("cp", filepath.Join(dir, "protoc/bin/protoc"), "/usr/local/bin/protoc"); err != nil {
			return err
		}
		fmt.Printf("install protoc %s success\n", version)
		return nil
	},
}
