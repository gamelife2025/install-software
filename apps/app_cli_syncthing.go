package apps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(syncthingApp)
}

func syncthingURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://github.com/syncthing/syncthing/releases/download/v%s/syncthing-linux-amd64-v%s.tar.gz", v, v), true
}

var syncthingApp = core.App{
	Name:     "syncthing",
	Category: "cli",
	Latest: func() (string, string, error) {
		tag, err := core.GithubLatestTag("syncthing/syncthing")
		if err != nil {
			return "", "", err
		}
		url, _ := syncthingURL(tag)
		return tag, url, nil
	},
	PinURL: syncthingURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "syncthing.tar.gz")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.ExtractTarUser(file, dir); err != nil {
			return err
		}
		bin := filepath.Join(dir, fmt.Sprintf("syncthing-linux-amd64-v%s/syncthing", core.TrimV(version)))
		if err := core.SudoRun("cp", bin, "/usr/local/bin/syncthing"); err != nil {
			return err
		}
		fmt.Printf("install syncthing %s success\n", version)
		return nil
	},
}
