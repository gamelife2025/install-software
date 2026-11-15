package apps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(javaApp)
}

var javaApp = core.App{
	Name:     "java",
	Category: "cli",
	Latest: func() (string, string, error) {
		return "8u231", "https://github.com/qilook/storage/releases/download/jdk/jdk-8u231-linux-x64.tar.gz", nil
	},
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "jdk-8u231-linux-x64.tar.gz")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.ExtractTarUser(file, dir); err != nil {
			return err
		}
		if err := core.SudoRun("rm", "-rf", "/opt/jdk"); err != nil {
			return err
		}
		if err := core.SudoRun("cp", "-r", filepath.Join(dir, "jdk1.8.0_231"), "/opt/jdk"); err != nil {
			return err
		}
		if err := core.WritePrivilegedFile("export PATH=$PATH:/opt/jdk/bin\n", "/etc/profile.d/jdk.sh", "+x"); err != nil {
			return err
		}
		fmt.Printf(`--------------------------------
install jdk1.8.0_231 success
you may need set:  export PATH=$PATH:/opt/jdk/bin
in your ~/.bashrc  or ~/.zshrc or /etc/zsh/zshenv
--------------------------------
`)
		return nil
	},
}
