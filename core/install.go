package core

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func haveSudo() bool {
	_, err := exec.LookPath("sudo")
	return err == nil
}

// sudoCmd 非 root 且有 sudo 时，用 sudo 执行
func sudoCmd(name string, args ...string) *exec.Cmd {
	if os.Geteuid() != 0 && haveSudo() {
		return exec.Command("sudo", append([]string{name}, args...)...)
	}
	return exec.Command(name, args...)
}

// Run 执行命令并挂接标准输出
func Run(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SudoRun 以 root 权限执行命令
func SudoRun(name string, args ...string) error {
	return Run(sudoCmd(name, args...))
}

type progressWriter struct {
	total, written, lastPct int64
}

func (w *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	w.written += int64(n)
	if w.total > 0 {
		pct := w.written * 100 / w.total
		if pct >= w.lastPct+10 || w.written == w.total {
			w.lastPct = pct
			fmt.Printf("\r下载进度: %d%%", pct)
		}
	} else {
		fmt.Printf("\r已下载: %d MB", w.written/1024/1024)
	}
	return n, nil
}

// Download 下载 url 到 dest（dest 应位于当前用户可写的临时目录）
func Download(url, dest string) error {
	resp, err := HTTPGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("下载失败 %s: %s", url, resp.Status)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	pw := &progressWriter{total: resp.ContentLength}
	if _, err := io.Copy(f, io.TeeReader(resp.Body, pw)); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

// FetchToTemp 下载 url 到临时目录，返回文件路径和所在目录（调用方负责 os.RemoveAll(dir)）
func FetchToTemp(url, name string) (file, dir string, err error) {
	dir, err = os.MkdirTemp("", "install-software-")
	if err != nil {
		return "", "", err
	}
	file = filepath.Join(dir, name)
	if err = Download(url, file); err != nil {
		os.RemoveAll(dir)
		return "", "", err
	}
	return file, dir, nil
}

// ExtractTarUser 以当前用户解压到 dir（GNU tar 自动识别 gz/xz）
func ExtractTarUser(archive, dir string) error {
	return Run(exec.Command("tar", "-xf", archive, "-C", dir))
}

// ExtractTarRoot 以 root 解压到 dir（/opt 等系统目录）
func ExtractTarRoot(archive, dir string, stripComponents int) error {
	args := []string{"--no-same-owner", "-xf", archive, "-C", dir}
	if stripComponents > 0 {
		args = append(args, "--strip-components", strconv.Itoa(stripComponents))
	}
	return SudoRun("tar", args...)
}

// InstallAppImage 安装 AppImage 到 /opt/<appDir>，并安装图标和 desktop 文件
func InstallAppImage(file, appDir, icon, desktop, desktopAs string) error {
	dest := "/opt/" + appDir
	if err := SudoRun("mkdir", "-p", dest); err != nil {
		return err
	}
	if err := SudoRun("cp", file, dest+"/"); err != nil {
		return err
	}
	if err := SudoRun("chmod", "+x", dest+"/"+filepath.Base(file)); err != nil {
		return err
	}
	if icon != "" {
		if err := InstallIcon(icon, dest); err != nil {
			return err
		}
	}
	if desktop != "" {
		if err := InstallDesktopFile(desktop, desktopAs); err != nil {
			return err
		}
	}
	return nil
}

// WritePrivilegedFile 写入需要 root 权限的文件（先写临时文件再 sudo cp）
func WritePrivilegedFile(content, dest string, mode string) error {
	tmp, err := os.CreateTemp("", "install-software-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(content); err != nil {
		tmp.Close()
		return err
	}
	tmp.Close()
	if err := SudoRun("cp", tmp.Name(), dest); err != nil {
		return err
	}
	return SudoRun("chmod", mode, dest)
}
