package apps

import (
	"fmt"
	"os"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(drawioApp)
}

func drawioURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://github.com/jgraph/drawio-desktop/releases/download/v%s/drawio-x86_64-%s.AppImage", v, v), true
}

var drawioApp = core.App{
	Name:     "drawio",
	Alias:    "2",
	Category: "desktop",
	Latest: func() (string, string, error) {
		tag, err := core.GithubLatestTag("jgraph/drawio-desktop")
		if err != nil {
			return "", "", err
		}
		url, _ := drawioURL(tag)
		return tag, url, nil
	},
	PinURL: drawioURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "drawio.AppImage")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.InstallAppImage(file, "drawio", "drawio.png", "drawio.desktop", ""); err != nil {
			return err
		}
		fmt.Printf("install drawio %s success\n", version)
		return nil
	},
}
