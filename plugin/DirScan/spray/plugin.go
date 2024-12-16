// ehole-------------------------------------
// @file      : plugin.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/11/21 21:55
// -------------------------------------------

package plugin

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/global"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"os"
	"path/filepath"
	"runtime"
)

func GetName() string {
	return "spray"
}

func Install() error {
	// 安装 Findomain 到指定目录下面
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
				fmt.Sprintf("Failed to set permissions: %v", err)
				return nil
			}
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

func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	data, ok := input.(types.AssetHttp)
	if !ok {
		return nil, nil
	}
	op.Log(fmt.Sprintf("scan terget begin: %v", data.URL))

	return nil, nil
}
