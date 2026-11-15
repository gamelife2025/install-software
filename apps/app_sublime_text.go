package apps

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(sublimeTextApp)
}

var sublimeBuildRe = regexp.MustCompile(`sublime_text_build_(\d+)_x64`)

func sublimeURL(build string) (string, bool) {
	b := core.TrimV(build)
	if b == "" {
		return "", false
	}
	return fmt.Sprintf("https://download.sublimetext.com/sublime_text_build_%s_x64.tar.xz", b), true
}

var sublimeTextApp = core.App{
	Name:     "sublime_text",
	Alias:    "4",
	Category: "desktop",
	Latest: func() (string, string, error) {
		resp, err := core.HTTPGet("https://www.sublimetext.com/download")
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", err
		}
		m := sublimeBuildRe.FindSubmatch(body)
		if m == nil {
			return "", "", fmt.Errorf("无法解析 sublime_text 最新 build")
		}
		url, _ := sublimeURL(string(m[1]))
		return string(m[1]), url, nil
	},
	PinURL: sublimeURL,
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, filepath.Base(url))
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.SudoRun("rm", "-rf", "/opt/sublime_text"); err != nil {
			return err
		}
		if err := core.ExtractTarRoot(file, "/opt", 0); err != nil {
			return err
		}
		if err := core.SudoRun("rm", "-f", "/usr/share/applications/sublime_text.desktop"); err != nil {
			return err
		}
		if err := core.SudoRun("ln", "-fs", "/opt/sublime_text/sublime_text.desktop", "/usr/share/applications/sublime_text.desktop"); err != nil {
			return err
		}
		if err := core.SudoRun("sed", "-i", `s/^Icon=.*/Icon=\/opt\/sublime_text\/Icon\/48x48\/sublime-text.png/`, "/opt/sublime_text/sublime_text.desktop"); err != nil {
			return err
		}
		wrapper := "#!/bin/sh\nexec /opt/sublime_text/sublime_text --fwdargv0 \"$0\" \"$@\"\n"
		if err := core.WritePrivilegedFile(wrapper, "/usr/local/bin/subl", "+x"); err != nil {
			return err
		}
		fmt.Println("install sublime_text latest success")
		return nil
	},
}
