// ehole-------------------------------------
// @file      : plugin.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/11/21 21:55
// -------------------------------------------

package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/global"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func GetName() string {
	return "EHole"
}

func Install() error {
	return nil
}

func Check() error {
	return nil
}

func Uninstall() error {
	return nil
}

type Fingerprints struct {
	Fingerprint []Fingerprint
}

type Fingerprint struct {
	Cms      string
	Method   string
	Location string
	Keyword  []string
}

func iskeyword(str string, keyword []string) bool {
	var x bool
	x = true
	for _, k := range keyword {
		if strings.Contains(str, k) {
			x = x && true
		} else {
			x = x && false
		}
	}
	return x
}

func isregular(str string, keyword []string) bool {
	var x bool
	x = true
	for _, k := range keyword {
		re := regexp.MustCompile(k)
		if re.Match([]byte(str)) {
			x = x && true
		} else {
			x = x && false
		}
	}
	return x
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	// 加载EHole指纹
	data, ok := input.(*types.AssetHttp)
	if !ok {
		// 说明不是http的资产，直接返回
		return nil, nil
	}
	op.Log(fmt.Sprintf("target %v begin", data.URL))
	parameter := op.Parameter
	var filgerFile string
	if parameter != "" {
		args, err := utils.Tools.ParseArgs(parameter, "finger")
		if err != nil {
		} else {
			for key, value := range args {
				if value != "" {
					switch key {
					case "finger":
						filgerFile = value
					}
				}
			}
		}
	}
	if filgerFile == "" {
		op.Log("EHole 指纹文件为空", "w")
		return nil, nil
	}
	fingerFilePath := filepath.Join(global.DictPath, filgerFile)
	content, err := os.ReadFile(fingerFilePath)
	if err != nil {
		op.Log(fmt.Sprintf("read finger error: %v", err))
		return nil, nil
	}
	var fingers Fingerprints
	err = json.Unmarshal(content, &fingers)
	if err != nil {
		op.Log(fmt.Sprintf("json to fingers error: %v", err))
		return nil, nil
	}
	// 使用 sync.Map 来保证并发安全
	uniqueCms := sync.Map{}
	for _, finger := range fingers.Fingerprint {
		select {
		case <-op.Ctx.Done():
			break
		default:
			if finger.Location == "body" {
				if finger.Method == "keyword" {
					if iskeyword(data.ResponseBody, finger.Keyword) {
						// 使用 sync.Map 进行并发安全操作
						uniqueCms.Store(finger.Cms, true)
					}
				}
				if finger.Method == "faviconhash" {
					if data.FavIconMMH3 == finger.Keyword[0] {
						uniqueCms.Store(finger.Cms, true)
					}
				}
				if finger.Method == "regular" {
					if isregular(data.ResponseBody, finger.Keyword) {
						uniqueCms.Store(finger.Cms, true)
					}
				}
			}

			if finger.Location == "header" {
				if finger.Method == "keyword" {
					if iskeyword(data.RawHeaders, finger.Keyword) {
						uniqueCms.Store(finger.Cms, true)
					}
				}
				if finger.Method == "regular" {
					if isregular(data.RawHeaders, finger.Keyword) {
						uniqueCms.Store(finger.Cms, true)
					}
				}
			}
			if finger.Location == "title" {
				if finger.Method == "keyword" {
					if iskeyword(data.Title, finger.Keyword) {
						uniqueCms.Store(finger.Cms, true)
					}
				}
				if finger.Method == "regular" {
					if isregular(data.Title, finger.Keyword) {
						uniqueCms.Store(finger.Cms, true)
					}
				}
			}
		}
	}
	// 从 sync.Map 中获取去重后的结果
	var result []string
	existingMap := make(map[string]bool) // 用于忽略大小写去重
	// 合并 data.Technologies 中的元素
	for _, v := range data.Technologies {
		lowerV := strings.ToLower(v) // 转换为小写进行去重
		if _, exists := existingMap[lowerV]; !exists {
			result = append(result, v) // 保持原始大小写
			existingMap[lowerV] = true // 标记该元素已添加
		}
	}

	// 将 uniqueCms 中的结果加入到 Technologies 中
	uniqueCms.Range(func(key, value interface{}) bool {
		cms := key.(string)
		lowerCms := strings.ToLower(cms) // 转换为小写进行去重
		if _, exists := existingMap[lowerCms]; !exists {
			result = append(result, cms) // 保持原始大小写
			existingMap[lowerCms] = true // 标记该元素已添加
		}
		return true
	})

	// 将去重后的结果赋值给 Technologies
	data.Technologies = result
	return nil, nil
}
