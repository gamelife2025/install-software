package apps

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(goApp)
}

func goURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://dl.google.com/go/go%s.linux-amd64.tar.gz", v), true
}

var goApp = core.App{
	Name:     "go",
	Category: "cli",
	Latest: func() (string, string, error) {
		resp, err := core.HTTPGet("https://go.dev/dl/?mode=json")
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return "", "", fmt.Errorf("go.dev: %s", resp.Status)
		}
		var releases []struct {
			Version string `json:"version"`
			Stable  bool   `json:"stable"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
			return "", "", err
		}
		for _, r := range releases {
			if r.Stable {
				url, _ := goURL(r.Version)
				return core.TrimV(r.Version), url, nil
			}
		}
		return "", "", fmt.Errorf("go.dev: 未找到稳定版本")
	},
	PinURL: goURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "go.tar.gz")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.SudoRun("rm", "-rf", "/usr/local/go"); err != nil {
			return err
		}
		if err := core.ExtractTarRoot(file, "/usr/local", 0); err != nil {
			return err
		}
		if err := core.SudoRun("ln", "-sf", "/usr/local/go/bin/go", "/usr/local/bin/go"); err != nil {
			return err
		}
		if err := core.SudoRun("ln", "-sf", "/usr/local/go/bin/gofmt", "/usr/local/bin/gofmt"); err != nil {
			return err
		}
		fmt.Printf("install go v%s success\n", core.TrimV(version))
		return nil
	},
}
