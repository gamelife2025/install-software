package apps

import (
	"fmt"
	"os"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(dockerComposeApp)
}

func dockerComposeURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://github.com/docker/compose/releases/download/v%s/docker-compose-Linux-x86_64", v), true
}

var dockerComposeApp = core.App{
	Name:     "docker-compose",
	Category: "cli",
	Latest: func() (string, string, error) {
		tag, err := core.GithubLatestTag("docker/compose")
		if err != nil {
			return "", "", err
		}
		url, _ := dockerComposeURL(tag)
		return tag, url, nil
	},
	PinURL: dockerComposeURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "docker-compose")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := os.Chmod(file, 0755); err != nil {
			return err
		}
		if err := core.SudoRun("cp", file, "/usr/local/bin/docker-compose"); err != nil {
			return err
		}
		fmt.Printf("install docker-compose %s success\n", version)
		return nil
	},
}
