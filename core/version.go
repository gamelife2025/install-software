package core

import (
	"regexp"
	"strconv"
	"strings"
)

var leadingNumRe = regexp.MustCompile(`^\d+`)

// versionParts 把版本号拆成数字段，非数字开头的段记为 0（避免 "2.0-beta" 之类 panic）
func versionParts(v string) []int {
	segs := strings.Split(strings.TrimLeft(v, "vV"), ".")
	parts := make([]int, 0, len(segs))
	for _, s := range segs {
		n := 0
		if m := leadingNumRe.FindString(s); m != "" {
			n, _ = strconv.Atoi(m)
		}
		parts = append(parts, n)
	}
	return parts
}

// VersionNewer 判断 latest 是否比 current 新（预发布 -rc 版本不算新）
func VersionNewer(current, latest string) bool {
	if strings.Contains(latest, "-rc") {
		return false
	}
	a, b := versionParts(current), versionParts(latest)
	for len(a) < len(b) {
		a = append(a, 0)
	}
	for len(b) < len(a) {
		b = append(b, 0)
	}
	for i := range a {
		if b[i] > a[i] {
			return true
		}
		if b[i] < a[i] {
			return false
		}
	}
	return false
}

// TrimV 去掉版本号前缀 v / go
func TrimV(v string) string {
	return strings.TrimPrefix(strings.TrimLeft(v, "vV"), "go")
}
