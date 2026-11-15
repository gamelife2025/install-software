package apps

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(pycharmApp)
}

type jetbrainsRelease struct {
	Version   string `json:"version"`
	Downloads struct {
		Linux struct {
			Link string `json:"link"`
		} `json:"linux"`
	} `json:"downloads"`
}

func pycharmURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	// 2025 起 JetBrains 将 Community/Professional 合并为统一发行版
	return fmt.Sprintf("https://download-cdn.jetbrains.com/python/pycharm-%s.tar.gz", v), true
}

var pycharmApp = core.App{
	Name:     "pycharm",
	Alias:    "6",
	Category: "desktop",
	Latest: func() (string, string, error) {
		resp, err := core.HTTPGet("https://data.services.jetbrains.com/products/releases?code=PCP&latest=true&type=release")
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return "", "", fmt.Errorf("jetbrains api: %s", resp.Status)
		}
		var data map[string][]jetbrainsRelease
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return "", "", err
		}
		releases := data["PCP"]
		if len(releases) == 0 {
			return "", "", fmt.Errorf("jetbrains api: 未找到 pycharm 版本信息")
		}
		return releases[0].Version, releases[0].Downloads.Linux.Link, nil
	},
	PinURL: pycharmURL,
	Install: func(version, url string) error {
		// 解压出的目录名与 tarball 文件名一致（pycharm-x.y 或旧版 pycharm-community-x.y）
		tarName := filepath.Base(url)
		extractedDir := strings.TrimSuffix(tarName, ".tar.gz")
		file, dir, err := core.FetchToTemp(url, tarName)
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.ExtractTarUser(file, dir); err != nil {
			return err
		}
		if err := core.SudoRun("rm", "-rf", "/opt/pycharm-community"); err != nil {
			return err
		}
		if err := core.SudoRun("mkdir", "-p", "/opt/pycharm-community"); err != nil {
			return err
		}
		if err := core.SudoRun("cp", "-R", filepath.Join(dir, extractedDir)+"/.", "/opt/pycharm-community/"); err != nil {
			return err
		}
		if err := core.InstallDesktopFile("pycharm.desktop", "jetbrains-pycharm-ce.desktop"); err != nil {
			return err
		}
		fmt.Printf("install pycharm-ce %s success\n", version)
		return nil
	},
}
