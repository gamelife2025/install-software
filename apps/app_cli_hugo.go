package apps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(hugoApp)
}

func hugoURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://github.com/gohugoio/hugo/releases/download/v%s/hugo_%s_Linux-64bit.tar.gz", v, v), true
}

var hugoApp = core.App{
	Name:     "hugo",
	Category: "cli",
	Latest: func() (string, string, error) {
		tag, err := core.GithubLatestTag("gohugoio/hugo")
		if err != nil {
			return "", "", err
		}
		url, _ := hugoURL(tag)
		return tag, url, nil
	},
	PinURL: hugoURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "hugo.tar.gz")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.ExtractTarUser(file, dir); err != nil {
			return err
		}
		if err := core.SudoRun("cp", filepath.Join(dir, "hugo"), "/usr/local/bin/hugo"); err != nil {
			return err
		}
		fmt.Printf("install hugo %s success\n", version)
		return nil
	},
}
