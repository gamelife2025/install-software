package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	// 注册所有 app（各文件 init() 自动注册）
	_ "github.com/gamelife2025/install-software/apps"
	"github.com/gamelife2025/install-software/core"
)

// 构建时通过 -ldflags "-X main.version=vX.Y.Z" 注入
var version = "dev"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "help", "-h", "--help":
		printHelp()
	case "list":
		printList()
	case "version", "-v", "--version":
		fmt.Println(version)
	case "upgrade":
		if err := selfUpgrade(); err != nil {
			fmt.Fprintln(os.Stderr, "升级失败:", err)
			os.Exit(1)
		}
	default:
		app := core.Find(args[0])
		if app == nil {
			fmt.Fprintf(os.Stderr, "未知命令: %s\n\n", args[0])
			printHelp()
			os.Exit(1)
		}
		pin := ""
		if len(args) > 1 {
			pin = args[1]
		}
		if err := installApp(app, pin); err != nil {
			fmt.Fprintf(os.Stderr, "安装 %s 失败: %v\n", app.Name, err)
			os.Exit(1)
		}
	}
}

func printHelp() {
	me := "install-software"
	fmt.Printf(`安装软件，一条命令搞定。安装时自动解析并下载最新版本。

用法:
  %s <app> [版本]     安装软件（版本可选，部分应用支持指定版本）
  %s list             列出所有支持的软件
  %s upgrade          升级 %s 自身到最新版本
  %s version          显示当前版本

`, me, me, me, me, me)
	printList()
}

func printList() {
	desktop := make([]core.App, 0)
	cli := make([]core.App, 0)
	for _, app := range core.List() {
		if app.Category == "desktop" {
			desktop = append(desktop, app)
		} else {
			cli = append(cli, app)
		}
	}
	// 各 app 文件的 init 注册顺序按文件名排序，展示时按数字别名排序
	sort.Slice(desktop, func(i, j int) bool {
		ai, _ := strconv.Atoi(desktop[i].Alias)
		aj, _ := strconv.Atoi(desktop[j].Alias)
		return ai < aj
	})
	sort.Slice(cli, func(i, j int) bool { return cli[i].Name < cli[j].Name })

	fmt.Println("------------------desktop-------------")
	for _, app := range desktop {
		fmt.Printf("  %-15s (install-software %s)\n", app.Name, app.Alias)
	}
	fmt.Println("------------------cli-----------------")
	for _, app := range cli {
		fmt.Printf("  %s\n", app.Name)
	}
}

// installApp 解析版本（或用用户指定的版本）并执行安装
func installApp(app *core.App, pin string) error {
	var ver, url string
	if pin != "" {
		if app.PinURL == nil {
			return fmt.Errorf("%s 不支持指定版本，请直接运行 install-software %s", app.Name, app.Name)
		}
		u, ok := app.PinURL(pin)
		if !ok {
			return fmt.Errorf("%s 不支持指定版本 %s", app.Name, pin)
		}
		ver, url = pin, u
	} else {
		fmt.Printf("正在解析 %s 最新版本...\n", app.Name)
		v, u, err := app.Latest()
		if err != nil {
			return err
		}
		ver, url = v, u
	}
	fmt.Printf("安装 %s %s\n%s\n", app.Name, ver, url)
	return app.Install(ver, url)
}
