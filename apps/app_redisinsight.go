package apps

import (
	"fmt"
	"os"

	"github.com/gamelife2025/install-software/core"
)

func init() {
	core.Register(redisInsightApp)
}

var redisInsightApp = core.App{
	Name:     "redisInsight",
	Alias:    "10",
	Category: "desktop",
	Latest: func() (string, string, error) {
		return "latest", "https://download.redisinsight.redis.com/latest/RedisInsight-linux-x86_64.AppImage", nil
	},
	Install: func(version, url string) error {
		file, dir, err := core.FetchToTemp(url, "RedisInsight.AppImage")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		if err := core.InstallAppImage(file, "redisInsight", "RedisInsight.png", "redisInsight.desktop", ""); err != nil {
			return err
		}
		fmt.Println("install redisInsight success")
		return nil
	},
}
