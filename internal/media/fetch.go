// via54Design — 设计模板引擎 + 叙事引擎
// Copyright (C) 2026  via54 (veawho)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// SPDX-License-Identifier: AGPL-3.0-only

package media

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const wikimediaAPI = "https://commons.wikimedia.org/w/api.php"
const userAgent = "via54Design/0.2 (https://github.com/veawho/via54Design; bot)"

// FetchResult 取图结果
type FetchResult struct {
	Path    string
	License string
	Author  string
	URL     string
}

// FetchImages 从 Wikimedia Commons 取图
func FetchImages(queries []string, outDir string, count int) ([]FetchResult, error) {
	os.MkdirAll(outDir, 0755)

	// 清代理 (Wikimedia 对代理 TLS 敏感)
	for _, k := range []string{"ALL_PROXY", "all_proxy", "HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy"} {
		os.Unsetenv(k)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	var results []FetchResult
	for _, q := range queries {
		params := url.Values{
			"action":       {"query"},
			"format":       {"json"},
			"generator":    {"search"},
			"gsrsearch":    {q},
			"gsrnamespace": {"6"},
			"gsrlimit":     {fmt.Sprintf("%d", count)},
			"prop":         {"imageinfo"},
			"iiprop":       {"url|extmetadata"},
			"iiurlwidth":   {"1200"},
		}
		req, _ := http.NewRequest("GET", wikimediaAPI+"?"+params.Encode(), nil)
		req.Header.Set("User-Agent", userAgent)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠️ 查询 '%s' 失败: %v\n", q, err)
			continue
		}
		var data struct {
			Query *struct {
				Pages map[string]struct {
					Title     string `json:"title"`
					ImageInfo []struct {
						ThumbURL string `json:"thumburl"`
						URL      string `json:"url"`
						ExtMeta  map[string]struct {
							Value string `json:"value"`
						} `json:"extmetadata"`
					} `json:"imageinfo"`
				} `json:"pages"`
			} `json:"query"`
		}
		json.NewDecoder(resp.Body).Decode(&data)
		resp.Body.Close()

		if data.Query == nil {
			continue
		}

		for _, page := range data.Query.Pages {
			if len(page.ImageInfo) == 0 {
				continue
			}
			ii := page.ImageInfo[0]
			thumb := ii.ThumbURL
			if thumb == "" {
				thumb = ii.URL
			}
			if thumb == "" {
				continue
			}

			// 下载
			ext := filepath.Ext(thumb)
			if ext == "" {
				ext = ".jpg"
			}
			filename := sanitize(q) + "_" + sanitize(strings.TrimPrefix(page.Title, "File:"))
			filename = truncate(filename, 55) + ext
			outPath := filepath.Join(outDir, filename)

			dlReq, _ := http.NewRequest("GET", thumb, nil)
			dlReq.Header.Set("User-Agent", userAgent)
			dlResp, err := client.Do(dlReq)
			if err != nil {
				continue
			}
			defer dlResp.Body.Close()

			f, _ := os.Create(outPath)
			written, _ := io.Copy(f, dlResp.Body)
			f.Close()

			if written > 1000 {
				license := ""
				if em, ok := ii.ExtMeta["LicenseShortName"]; ok {
					license = em.Value
				}
				author := ""
				if em, ok := ii.ExtMeta["Artist"]; ok {
					re := regexp.MustCompile("<[^>]+>")
					author = strings.TrimSpace(re.ReplaceAllString(em.Value, ""))
				}
				results = append(results, FetchResult{
					Path: outPath, License: license,
					Author: truncate(author, 60),
					URL:    ii.URL,
				})
				fmt.Printf("  ✅ %s  | %s\n", filepath.Base(outPath), license)
			}
		}
	}
	return results, nil
}

func sanitize(s string) string {
	re := regexp.MustCompile(`[^\w\-.]+`)
	return re.ReplaceAllString(s, "_")
}
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
