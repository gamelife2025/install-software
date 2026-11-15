package apps

import (
	"fmt"
	"os"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(etcdApp)
}

func etcdURL(version string) (string, bool) {
	v := core.TrimV(version)
	if v == "" {
		return "", false
	}
	return fmt.Sprintf("https://storage.googleapis.com/etcd/v%s/etcd-v%s-linux-amd64.tar.gz", v, v), true
}

var etcdApp = core.App{
	Name:     "etcd",
	Category: "cli",
	Latest: func() (string, string, error) {
		tag, err := core.GithubLatestTag("etcd-io/etcd")
		if err != nil {
			return "", "", err
		}
		url, _ := etcdURL(tag)
		return tag, url, nil
	},
	PinURL: etcdURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "etcd.tar.gz")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.SudoRun("rm", "-rf", "/opt/etcd"); err != nil {
			return err
		}
		if err := core.SudoRun("mkdir", "-p", "/opt/etcd"); err != nil {
			return err
		}
		if err := core.ExtractTarRoot(file, "/opt/etcd", 1); err != nil {
			return err
		}
		for _, bin := range []string{"etcdctl", "etcd", "etcdutl"} {
			if err := core.SudoRun("ln", "-sf", "/opt/etcd/"+bin, "/usr/local/bin/"+bin); err != nil {
				return err
			}
		}
		fmt.Printf("install etcd %s success\n", version)
		return nil
	},
}
