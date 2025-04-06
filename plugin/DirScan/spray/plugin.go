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
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

func GetName() string {
	return "spray"
}

func Install() error {
	sprayPath := filepath.Join(global.ExtDir, "spray")
	if err := os.MkdirAll(sprayPath, os.ModePerm); err != nil {
		logger.SlogError(fmt.Sprintf("Failed to create ksubdomain folder:", err))
		return err
	}
	osType := runtime.GOOS
	var sprayURL string
	var fileName string
	switch osType {
	case "windows":
		sprayURL = "https://github.com/chainreactors/spray/releases/download/v1.1.6/spray_windows_amd64.exe"
		fileName = "spray.exe"
	case "linux":
		sprayURL = "https://github.com/chainreactors/spray/releases/download/v1.1.6/spray_linux_amd64"
		fileName = "spray"
	}
	sprayExecPath := filepath.Join(sprayPath, fileName)
	if _, err := os.Stat(sprayExecPath); os.IsNotExist(err) {
		downloadPath := filepath.Join(sprayPath, fileName) // 临时下载路径
		success, err := utils.Tools.HttpGetDownloadFile(sprayURL, downloadPath)
		if err != nil || !success {
			logger.SlogErrorLocal(fmt.Sprintf("Failed to download spray: %v", err))
			return err
		}

		logger.SlogInfo("spray Download successful")
		switch osType {
		case "linux":
			err = os.Chmod(sprayExecPath, 0755)
			if err != nil {
				logger.SlogErrorLocal(fmt.Sprintf("Failed to set permissions: %v", err))
				return nil
			}
		}
		err = utils.Tools.EnsureDir(filepath.Join(sprayPath, "result"))
		if err != nil {
			logger.SlogErrorLocal(fmt.Sprintf("Failed to EnsureDir: %v", err))
			return err
		}
		logger.SlogInfo("spray installed successfully")
	}
	return nil
}

func Check() error {
	return nil
}

func Uninstall() error {
	return nil
}

type Framework struct {
	Name string `json:"name"`
}

type Frameworks map[string]*Framework

type SprayResult struct {
	Number       int        `json:"number"`
	Parent       int        `json:"parent"`
	IsValid      bool       `json:"valid"`
	IsFuzzy      bool       `json:"fuzzy"`
	UrlString    string     `json:"url"`
	Path         string     `json:"path"`
	Host         string     `json:"host"`
	BodyLength   int        `json:"body_length"`
	ExceedLength bool       `json:"-"`
	HeaderLength int        `json:"header_length"`
	RedirectURL  string     `json:"redirect_url,omitempty"`
	FrontURL     string     `json:"front_url,omitempty"`
	Status       int        `json:"status"`
	Spended      int64      `json:"spend"` // 耗时, 毫秒
	ContentType  string     `json:"content_type"`
	Title        string     `json:"title"`
	Frameworks   Frameworks `json:"frameworks"`
	ErrString    string     `json:"error"`
	Reason       string     `json:"reason"`
	ReqDepth     int        `json:"depth"`
	Distance     uint8      `json:"distance"`
	Unique       uint16     `json:"unique"`
}

func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	data, ok := input.(types.AssetHttp)
	if !ok {
		return nil, nil
	}
	op.Log(fmt.Sprintf("scan terget begin: %v", data.URL))
	start := time.Now()
	parameter := op.Parameter
	executionTimeout := 20
	var cmdArray []string
	cmdArray = append(cmdArray, "-u")
	cmdArray = append(cmdArray, data.URL)
	// outputDirName := fmt.Sprintf("bbot_result_%s", targetMD5)
	if parameter != "" {
		args, err := utils.Tools.ParseArgs(parameter, "et", "finger", "common", "bak", "d")
		if err != nil {
			op.Log(fmt.Sprintf("parse args error: %v", err), "e")
		} else {
			for key, value := range args {
				if value != "" {
					switch key {
					case "d":
						cmdArray = append(cmdArray, "-d")
						dirDicConfigPath := filepath.Join(global.DictPath, value)
						cmdArray = append(cmdArray, dirDicConfigPath)
					case "finger":
						if value == "true" {
							cmdArray = append(cmdArray, "--finger")
						}
					case "common":
						if value == "true" {
							cmdArray = append(cmdArray, "--common")
						}
					case "bak":
						if value == "true" {
							cmdArray = append(cmdArray, "--bak")
						}
					case "et":
						executionTimeout, _ = strconv.Atoi(value)
					default:
						continue
					}
				}

			}
		}
	}
	cmdArray = append(cmdArray, "--json")
	cmdArray = append(cmdArray, "-f")
	fileName := utils.Tools.GenerateRandomString(6) + ".json"
	resultFile := filepath.Join(global.ExtDir, "spray", "result", fileName)
	cmdArray = append(cmdArray, resultFile)
	//defer utils.Tools.DeleteFile(resultFile)
	osType := runtime.GOOS
	var cmd string
	switch osType {
	case "windows":
		cmd = "spray.exe"
	default:
		cmd = "spray"
	}

	err := utils.Tools.ExecuteCommandWithTimeout(filepath.Join(global.ExtDir, "spray", cmd), cmdArray, time.Duration(executionTimeout)*time.Minute, op.Ctx)
	if err != nil {
		op.Log(fmt.Sprintf("ExecuteCommandWithTimeout error: %v", err), "e")
		return nil, err
	}
	resultsChan := make(chan string, 10)
	go func() {
		err := utils.Tools.ReadFileLineReader(resultFile, resultsChan, op.Ctx)
		if err != nil {
			op.Log(fmt.Sprintf("target %v read result file error: %v", data, err), "e")
		}
	}()
	for line := range resultsChan {
		var result SprayResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			op.Log(fmt.Sprintf("json.Unmarshal %v error: %v", line, err))
			continue
		}
		if result.UrlString == data.URL {
			continue
		}
		var DirResult types.DirResult
		DirResult.Url = result.UrlString
		DirResult.Status = result.Status
		DirResult.Length = result.BodyLength
		var tags []string
		if result.Title != "" {
			tags = append(tags, "title:"+result.Title)
		}
		for _, framework := range result.Frameworks {
			if framework.Name != "" {
				tags = append(tags, "app:"+framework.Name)
			}
		}
		if result.ContentType != "" {
			tags = append(tags, "content:"+result.ContentType)
		}
		DirResult.Tags = tags
		op.ResultFunc(DirResult)
	}
	duration := time.Since(start)
	op.Log(fmt.Sprintf("spray scan completed: target %v time: %v", data.URL, duration))
	return nil, nil
}
