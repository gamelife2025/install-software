package core

// App 一个可安装的软件
type App struct {
	Name     string // 命令名
	Alias    string // 数字别名，无则为 ""
	Category string // "desktop" 或 "cli"
	// Latest 解析最新版本，返回版本号和下载地址
	Latest func() (version, url string, err error)
	// PinURL 由用户指定的版本构造下载地址；nil 表示不支持指定版本
	PinURL func(version string) (url string, ok bool)
	// Install 下载并安装指定版本
	Install func(version, url string) error
}

var apps []App

// Register 注册应用，由各 app 文件的 init() 调用
func Register(list ...App) {
	apps = append(apps, list...)
}

// List 返回所有已注册应用
func List() []App {
	return apps
}

// Find 按命令名或数字别名查找应用
func Find(name string) *App {
	for i := range apps {
		if apps[i].Name == name || (apps[i].Alias != "" && apps[i].Alias == name) {
			return &apps[i]
		}
	}
	return nil
}
