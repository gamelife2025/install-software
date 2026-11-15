package apps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(vscodeCLIApp)
}

var vscodeCLIApp = core.App{
	Name:     "vscode-cli",
	Category: "cli",
	Latest: func() (string, string, error) {
		commitID, _, err := resolveVSCode()
		if err != nil {
			return "", "", err
		}
		return commitID, fmt.Sprintf("https://az764295.vo.msecnd.net/stable/%s/vscode_cli_alpine_x64_cli.tar.gz", commitID), nil
	},
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "vscode_cli.tar.gz")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.ExtractTarUser(file, dir); err != nil {
			return err
		}
		if err := core.SudoRun("mv", filepath.Join(dir, "code"), "/usr/local/bin/vscode-cli"); err != nil {
			return err
		}
		fmt.Println("install vscode-cli success")
		return nil
	},
}
