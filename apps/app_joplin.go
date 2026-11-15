package apps

import (
	"fmt"
	"os"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(joplinApp)
}

func joplinURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://github.com/laurent22/joplin/releases/download/v%s/Joplin-%s.AppImage", v, v), true
}

var joplinApp = core.App{
	Name:     "joplin",
	Alias:    "7",
	Category: "desktop",
	Latest: func() (string, string, error) {
		tag, err := core.GithubLatestTag("laurent22/joplin")
		if err != nil {
			return "", "", err
		}
		url, _ := joplinURL(tag)
		return tag, url, nil
	},
	PinURL: joplinURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "Joplin.AppImage")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.InstallAppImage(file, "joplin", "joplin.png", "joplin.desktop", "joplinapp.desktop"); err != nil {
			return err
		}
		fmt.Printf("install joplin %s success\n", version)
		return nil
	},
}
