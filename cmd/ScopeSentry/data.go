// main-------------------------------------
// @file      : data.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/11/16 18:01
// -------------------------------------------

package main

import (
	"encoding/json"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
)

var (
	AssetHttpData = types.AssetHttp{
		Time:          "2024-11-16T08:30:00Z",
		LastScanTime:  "2024-11-15T10:00:00Z",
		TLSData:       nil, // 如果没有 TLS 数据可以设置为 nil
		Hashes:        map[string]interface{}{"sha256": "abcdef1234567890"},
		CDNName:       "Cloudflare",
		Port:          "443",
		URL:           "https://example.com",
		Title:         "Example Site",
		Type:          "web",
		Error:         "",
		ResponseBody:  "Example response body here.",
		Host:          "example.com",
		IP:            "192.168.1.1",
		Screenshot:    "path/to/screenshot.png",
		FavIconMMH3:   "abc123",
		FaviconPath:   "/assets/favicon.ico",
		RawHeaders:    "HTTP/1.1 200 OK\nContent-Type: text/html\n",
		Jarm:          "JARM hash data here",
		Technologies:  []string{"NGINX", "React", "Node.js"},
		StatusCode:    200,
		ContentLength: 1234,
		CDN:           true,
		Webcheck:      true,
		Project:       "Project X",
		IconContent:   "base64iconcontenthere",
		Domain:        "example.com",
		TaskName:      []string{"Task 1", "Task 2"},
		WebServer:     "nginx/1.21.0",
		Service:       "Web Hosting",
		RootDomain:    "example.com",
		Tags:          []string{"production", "ssl", "security"},
	}
	AssetOtherData = types.AssetOther{
		Time:         "2024-11-16T08:30:00Z",
		LastScanTime: "2024-11-15T10:00:00Z",
		Host:         "other-example.com",
		IP:           "192.168.1.2",
		Port:         "8080",
		Service:      "HTTP Server",
		TLS:          true,
		Transport:    "TCP",
		Version:      "1.0.0",
		Raw:          json.RawMessage(`{"metadataField": "value"}`),
		Project:      "Project Y",
		Type:         "service",
		Tags:         []string{"internal", "ssl", "test"},
		TaskName:     []string{"Task A", "Task B"},
		RootDomain:   "other-example.com",
	}
)
