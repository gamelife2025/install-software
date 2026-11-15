package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 60 * time.Second}

// httpNoRedirect 不跟随重定向，用于解析 firefox/vscode 的真实下载地址
var httpNoRedirect = &http.Client{
	Timeout: 60 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// HTTPGet 带重试的 GET 请求
func HTTPGet(url string) (*http.Response, error) {
	return getWithRetry(httpClient, url)
}

// HTTPGetNoRedirect 带重试的 GET 请求，不跟随重定向
func HTTPGetNoRedirect(url string) (*http.Response, error) {
	return getWithRetry(httpNoRedirect, url)
}

func getWithRetry(client *http.Client, url string) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "install-software")
		if token := os.Getenv("GITHUB_TOKEN"); token != "" && strings.Contains(url, "api.github.com") {
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Accept", "application/vnd.github+json")
		}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < 2 {
				fmt.Printf("请求失败，正在重试 %d/2: %s (%v)\n", attempt+1, url, err)
				time.Sleep(time.Second)
				continue
			}
			return nil, err
		}
		return resp, nil
	}
	return nil, lastErr
}

type GHAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type GHRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []GHAsset `json:"assets"`
}

// GithubLatestRelease 获取 GitHub 仓库最新 release
func GithubLatestRelease(repo string) (*GHRelease, error) {
	resp, err := HTTPGet("https://api.github.com/repos/" + repo + "/releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api %s: %s", repo, resp.Status)
	}
	var rel GHRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	if rel.TagName == "" {
		return nil, fmt.Errorf("github api %s: 未找到 tag", repo)
	}
	return &rel, nil
}

// GithubLatestTag 获取 GitHub 仓库最新 release 的 tag
func GithubLatestTag(repo string) (string, error) {
	rel, err := GithubLatestRelease(repo)
	if err != nil {
		return "", err
	}
	return rel.TagName, nil
}

// ResolveRedirect 获取重定向后的地址（用于 firefox/vscode 等）
func ResolveRedirect(url string) (string, error) {
	resp, err := HTTPGetNoRedirect(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", fmt.Errorf("%s 未返回重定向地址 (status %s)", url, resp.Status)
	}
	return loc, nil
}
