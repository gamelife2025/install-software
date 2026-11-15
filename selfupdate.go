package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gamelife2025/install-software/core"
)

const selfRepo = "gamelife2025/install-software"
const selfAsset = "install-software-linux-amd64"

// selfUpgrade 从 GitHub Releases 下载最新版本并替换当前二进制
func selfUpgrade() error {
	fmt.Println("当前版本:", version)
	rel, err := core.GithubLatestRelease(selfRepo)
	if err != nil {
		return err
	}
	fmt.Println("最新版本:", rel.TagName)
	if rel.TagName == version {
		fmt.Println("已是最新版本")
		return nil
	}

	var binaryURL, shaURL string
	for _, a := range rel.Assets {
		switch a.Name {
		case selfAsset:
			binaryURL = a.BrowserDownloadURL
		case selfAsset + ".sha256":
			shaURL = a.BrowserDownloadURL
		}
	}
	if binaryURL == "" {
		return fmt.Errorf("release %s 中未找到 %s", rel.TagName, selfAsset)
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return err
	}

	file, dir, err := core.FetchToTemp(binaryURL, selfAsset)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	// 校验 sha256（release 提供时）
	if shaURL != "" {
		shaFile := filepath.Join(dir, "sha256")
		if err := core.Download(shaURL, shaFile); err != nil {
			return fmt.Errorf("下载校验文件失败: %w", err)
		}
		if err := verifySHA256(file, shaFile); err != nil {
			return err
		}
		fmt.Println("sha256 校验通过")
	}

	if err := os.Chmod(file, 0755); err != nil {
		return err
	}

	// 先尝试直接替换（rename 失败则拷贝，权限不足则 sudo）
	if err := os.Rename(file, exe); err == nil {
		fmt.Printf("升级成功: %s -> %s\n", version, rel.TagName)
		return nil
	}
	if err := copyFile(file, exe); err == nil {
		fmt.Printf("升级成功: %s -> %s\n", version, rel.TagName)
		return nil
	}
	if err := core.SudoRun("cp", file, exe); err != nil {
		return err
	}
	if err := core.SudoRun("chmod", "+x", exe); err != nil {
		return err
	}
	fmt.Printf("升级成功: %s -> %s\n", version, rel.TagName)
	return nil
}

func verifySHA256(file, shaFile string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))

	data, err := os.ReadFile(shaFile)
	if err != nil {
		return err
	}
	// sha256sum 输出格式: "<hex>  <filename>"
	want := strings.Fields(string(data))
	if len(want) == 0 {
		return fmt.Errorf("校验文件内容为空")
	}
	if got != want[0] {
		return fmt.Errorf("sha256 校验失败: 期望 %s, 实际 %s", want[0], got)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
