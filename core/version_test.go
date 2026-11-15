package core

import "testing"

func TestVersionNewer(t *testing.T) {
	cases := []struct {
		current, latest string
		want            bool
	}{
		{"v1.2.3", "v1.3", true},
		{"v1.3.3", "v1.3.4", true},
		{"v1.3", "v1.3.4", true},
		{"v3.5.4", "v3.4.20", false},
		{"v17.9.1", "v16.17.0", false},
		{"v1.2.3", "v1.2.2.1", false},
		{"v1.2.3", "v1.2.3", false},
		// 预发布版本不算新
		{"v1.2.3", "v1.3.0-rc1", false},
		// 非数字段不应 panic
		{"v1.2.3", "v1.2.3-beta", false},
	}
	for _, c := range cases {
		if got := VersionNewer(c.current, c.latest); got != c.want {
			t.Errorf("VersionNewer(%q, %q) = %v, want %v", c.current, c.latest, got, c.want)
		}
	}
}

func TestTrimV(t *testing.T) {
	cases := map[string]string{
		"v1.2.3":   "1.2.3",
		"V1.2.3":   "1.2.3",
		"1.2.3":    "1.2.3",
		"go1.27.4": "1.27.4",
	}
	for in, want := range cases {
		if got := TrimV(in); got != want {
			t.Errorf("TrimV(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFind(t *testing.T) {
	Register(App{Name: "test-dummy", Alias: "99", Category: "cli"})
	if Find("test-dummy") == nil {
		t.Error("Find(test-dummy) = nil")
	}
	if Find("99") == nil {
		t.Error("Find(99) = nil，别名未命中")
	}
	if Find("not-exist") != nil {
		t.Error("Find(not-exist) 应为 nil")
	}
}
