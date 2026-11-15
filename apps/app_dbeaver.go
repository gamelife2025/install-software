package apps

import (
	"fmt"
	"os"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(dbeaverApp)
}

var dbeaverApp = core.App{
	Name:     "dbeaver",
	Alias:    "1",
	Category: "desktop",
	Latest: func() (string, string, error) {
		// 官方提供固定 latest 地址
		return "latest", "https://dbeaver.io/files/dbeaver-ce-latest-linux.gtk.x86_64.tar.gz", nil
	},
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "dbeaver-ce-latest-linux.gtk.x86_64.tar.gz")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.ExtractTarRoot(file, "/opt", 0); err != nil {
			return err
		}
		if err := core.InstallDesktopFile("dbeaver.desktop", ""); err != nil {
			return err
		}
		fmt.Println("install dbeaver latest success")
		return nil
	},
}
