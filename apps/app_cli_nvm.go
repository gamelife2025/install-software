package apps

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(nvmApp)
}

func nvmURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://raw.githubusercontent.com/nvm-sh/nvm/v%s/install.sh", v), true
}

var nvmApp = core.App{
	Name:     "nvm",
	Category: "cli",
	Latest: func() (string, string, error) {
		tag, err := core.GithubLatestTag("nvm-sh/nvm")
		if err != nil {
			return "", "", err
		}
		url, _ := nvmURL(tag)
		return tag, url, nil
	},
	PinURL: nvmURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "nvm-install.sh")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		// nvm 安装到用户 home，不要用 sudo
		if err := core.Run(exec.Command("bash", file)); err != nil {
			return err
		}
		fmt.Printf("install nvm %s success, 请重新打开终端或 source ~/.bashrc / ~/.zshrc\n", version)
		return nil
	},
}
