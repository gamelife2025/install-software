package apps

import (
	"os"
	"testing"

	"github.com/gamelife2025/install-software/core"
)

// 以下 Resolver 测试需要访问外网，设置 RUN_NETWORK_TESTS=1 后运行
func skipUnlessNetwork(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_NETWORK_TESTS") != "1" {
		t.Skip("需要网络，设置 RUN_NETWORK_TESTS=1 运行")
	}
}

func TestQQRegex(t *testing.T) {
	js := `var x = "https://dldir1.qq.com/qqfile/qq/QQNT/Linux/QQ_3.2.9_240606_x86_64_01.AppImage";`
	m := qqRe.FindStringSubmatch(js)
	if m == nil {
		t.Fatal("qq 正则未匹配")
	}
	if m[1] != "3.2.9" {
		t.Errorf("version = %q, want 3.2.9", m[1])
	}
}

func TestQQMusicRegex(t *testing.T) {
	html := `<a href="https://c.y.qq.com/cgi-bin/file_redirect.fcg?bid=dldir&file=ecosfile_plink%2Fmusic_clntupate%2Flinux%2Fother%2Fqqmusic-1.1.8.AppImage&sign=1-abc">下载</a>`
	m := qqmusicRe.FindStringSubmatch(html)
	if m == nil {
		t.Fatal("qqmusic 正则未匹配")
	}
	if m[1] != "1.1.8" {
		t.Errorf("version = %q, want 1.1.8", m[1])
	}
}

func TestFirefoxRegex(t *testing.T) {
	url := "https://download-installer.cdn.mozilla.net/pub/firefox/releases/160.0/linux-x86_64/zh-CN/firefox-160.0.tar.xz"
	m := firefoxRe.FindStringSubmatch(url)
	if m == nil {
		t.Fatal("firefox 正则未匹配")
	}
	if m[1] != "160.0" {
		t.Errorf("version = %q, want 160.0", m[1])
	}
}

func TestPinURLs(t *testing.T) {
	cases := []struct {
		name    string
		fn      func(string) (string, bool)
		version string
		wantURL string
		wantOK  bool
	}{
		{"go", goURL, "1.27.4", "https://dl.google.com/go/go1.27.4.linux-amd64.tar.gz", true},
		{"go-strip-v", goURL, "v1.27.4", "https://dl.google.com/go/go1.27.4.linux-amd64.tar.gz", true},
		{"gh", ghURL, "v2.104.0", "https://github.com/cli/cli/releases/download/v2.104.0/gh_2.104.0_linux_amd64.tar.gz", true},
		{"etcd", etcdURL, "v3.7.3", "https://storage.googleapis.com/etcd/v3.7.3/etcd-v3.7.3-linux-amd64.tar.gz", true},
		{"protoc", protocURL, "v36.2", "https://github.com/protocolbuffers/protobuf/releases/download/v36.2/protoc-36.2-linux-x86_64.zip", true},
		{"protoc-too-old", protocURL, "v2.5.0", "", false},
		{"upx", upxURL, "v5.3.0", "https://github.com/upx/upx/releases/download/v5.3.0/upx-5.3.0-amd64_linux.tar.xz", true},
		{"syncthing", syncthingURL, "v2.1.7", "https://github.com/syncthing/syncthing/releases/download/v2.1.7/syncthing-linux-amd64-v2.1.7.tar.gz", true},
		{"hugo", hugoURL, "v0.170.0", "https://github.com/gohugoio/hugo/releases/download/v0.170.0/hugo_0.170.0_Linux-64bit.tar.gz", true},
		{"drawio", drawioURL, "v32.6.10", "https://github.com/jgraph/drawio-desktop/releases/download/v32.6.10/drawio-x86_64-32.6.10.AppImage", true},
		{"joplin", joplinURL, "v3.7.22", "https://github.com/laurent22/joplin/releases/download/v3.7.22/Joplin-3.7.22.AppImage", true},
		{"docker-compose", dockerComposeURL, "v5.7.1", "https://github.com/docker/compose/releases/download/v5.7.1/docker-compose-Linux-x86_64", true},
		{"nvm", nvmURL, "v0.41.0", "https://raw.githubusercontent.com/nvm-sh/nvm/v0.41.0/install.sh", true},
		{"pycharm", pycharmURL, "2026.3", "https://download-cdn.jetbrains.com/python/pycharm-2026.3.tar.gz", true},
		{"sublime", sublimeURL, "4213", "https://download.sublimetext.com/sublime_text_build_4213_x64.tar.xz", true},
		{"firefox", firefoxURL, "160.0", "https://download-installer.cdn.mozilla.net/pub/firefox/releases/160.0/linux-x86_64/zh-CN/firefox-160.0.tar.xz", true},
		{"empty", goURL, "", "", false},
	}
	for _, c := range cases {
		url, ok := c.fn(c.version)
		if ok != c.wantOK || url != c.wantURL {
			t.Errorf("%s: PinURL(%q) = (%q, %v), want (%q, %v)", c.name, c.version, url, ok, c.wantURL, c.wantOK)
		}
	}
}

// TestAllAppsRegistered 确认 22 个 app 全部注册
func TestAllAppsRegistered(t *testing.T) {
	want := []string{
		"dbeaver", "drawio", "firefox", "sublime_text", "vscode",
		"pycharm", "joplin", "qq", "qqmusic", "redisInsight", "hoppscotch",
		"docker-compose", "etcd", "gh", "go", "hugo", "java",
		"protoc", "syncthing", "upx", "vscode-cli", "nvm",
	}
	registered := make(map[string]bool)
	for _, app := range core.List() {
		registered[app.Name] = true
	}
	for _, name := range want {
		if !registered[name] {
			t.Errorf("app %q 未注册", name)
		}
	}
}

func TestResolverGo(t *testing.T) {
	skipUnlessNetwork(t)
	v, url, err := goApp.Latest()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("go %s %s", v, url)
}

func TestResolverFirefox(t *testing.T) {
	skipUnlessNetwork(t)
	v, url, err := firefoxApp.Latest()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("firefox %s %s", v, url)
}

func TestResolverVSCode(t *testing.T) {
	skipUnlessNetwork(t)
	v, url, err := vscodeApp.Latest()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("vscode %s %s", v, url)
}

func TestResolverPycharm(t *testing.T) {
	skipUnlessNetwork(t)
	v, url, err := pycharmApp.Latest()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("pycharm %s %s", v, url)
}

func TestResolverQQ(t *testing.T) {
	skipUnlessNetwork(t)
	v, url, err := qqApp.Latest()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("qq %s %s", v, url)
}

func TestResolverQQMusic(t *testing.T) {
	skipUnlessNetwork(t)
	v, url, err := qqmusicApp.Latest()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("qqmusic %s %s", v, url)
}

func TestResolverSublime(t *testing.T) {
	skipUnlessNetwork(t)
	v, url, err := sublimeTextApp.Latest()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("sublime_text %s %s", v, url)
}
