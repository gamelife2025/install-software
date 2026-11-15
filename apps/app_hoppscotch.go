package apps

import (
	"fmt"
	"os"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(hoppscotchApp)
}

var hoppscotchApp = core.App{
	Name:     "hoppscotch",
	Alias:    "11",
	Category: "desktop",
	Latest: func() (string, string, error) {
		return "latest", "https://github.com/hoppscotch/releases/releases/latest/download/Hoppscotch_linux_x64.AppImage", nil
	},
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "Hoppscotch_linux_x64.AppImage")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.InstallAppImage(file, "hoppscotch", "hoppscotch.png", "hoppscotch.desktop", ""); err != nil {
			return err
		}
		fmt.Println("install hoppscotch success")
		return nil
	},
}
