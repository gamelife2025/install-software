package apps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(ghApp)
}

func ghURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://github.com/cli/cli/releases/download/v%s/gh_%s_linux_amd64.tar.gz", v, v), true
}

var ghApp = core.App{
	Name:     "gh",
	Category: "cli",
	Latest: func() (string, string, error) {
		tag, err := core.GithubLatestTag("cli/cli")
		if err != nil {
			return "", "", err
		}
		url, _ := ghURL(tag)
		return tag, url, nil
	},
	PinURL: ghURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "gh.tar.gz")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.ExtractTarUser(file, dir); err != nil {
			return err
		}
		bin := filepath.Join(dir, fmt.Sprintf("gh_%s_linux_amd64/bin/gh", core.TrimV(version)))
		if err := core.SudoRun("cp", bin, "/usr/local/bin/gh"); err != nil {
			return err
		}
		fmt.Printf("install gh %s success\n", version)
		return nil
	},
}
