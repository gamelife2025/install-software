package core

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed desktop icon
var assets embed.FS

// ReadAsset 读取内嵌资源内容
func ReadAsset(name string) ([]byte, error) {
	return assets.ReadFile(name)
}

// InstallDesktopFile 把内嵌的 desktop 文件安装到 /usr/share/applications/
// asName 为空时使用原名
func InstallDesktopFile(name, asName string) error {
	if asName == "" {
		asName = name
	}
	data, err := assets.ReadFile("desktop/" + name)
	if err != nil {
		return err
	}
	tmp := filepath.Join(os.TempDir(), "install-software-"+name)
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	defer os.Remove(tmp)
	return SudoRun("cp", tmp, "/usr/share/applications/"+asName)
}

// InstallIcon 把内嵌的图标拷贝到 destDir
func InstallIcon(name, destDir string) error {
	data, err := assets.ReadFile("icon/" + name)
	if err != nil {
		return err
	}
	tmp := filepath.Join(os.TempDir(), "install-software-"+name)
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	defer os.Remove(tmp)
	return SudoRun("cp", tmp, destDir+"/"+name)
}
