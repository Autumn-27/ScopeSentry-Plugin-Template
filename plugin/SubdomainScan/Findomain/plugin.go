package plugin

import (
	"archive/zip"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/global"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetName() string {
	return "Findomain"
}

func Install() error {
	// 安装 Findomain 到指定目录下面
	findomainPath := filepath.Join(global.ExtDir, "Findomain")
	if err := os.MkdirAll(findomainPath, os.ModePerm); err != nil {
		logger.SlogError(fmt.Sprintf("Failed to create ksubdomain folder:", err))
		return err
	}
	osType := runtime.GOOS
	var findomainURL string
	var fileName string
	switch osType {
	case "windows":
		findomainURL = "https://github.com/Findomain/Findomain/releases/download/9.0.4/findomain-windows.exe.zip"
		fileName = "findomain.exe"
	case "linux":
		findomainURL = "https://github.com/Findomain/Findomain/releases/download/9.0.4/findomain-linux.zip"
		fileName = "findomain"
	}
	findomainExecPath := filepath.Join(findomainPath, fileName)
	if _, err := os.Stat(findomainExecPath); os.IsNotExist(err) {
		downloadPath := filepath.Join(global.ExtDir, "Findomain", "findomain.zip") // 临时下载路径
		success, err := utils.Tools.HttpGetDownloadFile(findomainURL, downloadPath)
		if err != nil || !success {
			logger.SlogErrorLocal(fmt.Sprintf("Failed to download findomain: %v", err))
			return err
		}

		logger.SlogInfo("Findomain Download successful")

		err = Unzip(downloadPath, findomainPath)
		if err != nil {
			fmt.Printf("Failed to extract findomain: %v", err)
			return nil
		}
		switch osType {
		case "linux":
			findomainExecPath := filepath.Join(findomainPath, "findomain")
			err = os.Chmod(findomainExecPath, 0755)
			if err != nil {
				fmt.Sprintf("Failed to set permissions: %v", err)
				return nil
			}
		}
		defer utils.Tools.DeleteFile(downloadPath)
		logger.SlogInfo("Findomain installed successfully")
	}
	return nil
}

// Unzip 解压缩一个 ZIP 文件到目标目录
func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func Check() error {
	osType := runtime.GOOS
	var fileName string
	switch osType {
	case "windows":
		fileName = "findomain.exe"
	case "linux":
		fileName = "findomain"
	}
	findomainPath := filepath.Join(global.ExtDir, "Findomain")
	findomainExecPath := filepath.Join(findomainPath, fileName)
	// 判断文件是否存在
	if _, err := os.Stat(findomainExecPath); os.IsNotExist(err) {
		logger.SlogInfo("Findomain is not installed")
		return fmt.Errorf("findomain is not installed")
	}
	logger.SlogInfo("check over - Findomain is installed")
	return nil
}

func Uninstall() error {
	return nil
}

func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	// 根据目标生成子域名，将子域名发送到result
	// 	err := logger.NewLogger()
	// 记录 debug 级别日志
	// op.Log("test-a","d")
	data, ok := input.(string)
	if !ok {
		// 说明不是http的资产，直接返回
		return nil, nil
	}
	osType := runtime.GOOS
	var fileName string
	switch osType {
	case "windows":
		fileName = "findomain.exe"
	case "linux":
		fileName = "findomain"
	}
	op.Log(fmt.Sprintf("target %v begin scan", data))
	start := time.Now()
	findomainPath := filepath.Join(global.ExtDir, "Findomain")
	findomainExecPath := filepath.Join(findomainPath, fileName)
	parameter := op.Parameter
	executionTimeout := 20
	if parameter != "" {
		args, err := utils.Tools.ParseArgs(parameter, "et")
		if err != nil {
		} else {
			for key, value := range args {
				if value != "" {
					switch key {
					case "et":
						executionTimeout, _ = strconv.Atoi(value)
					}
				}
			}
		}
	}
	resultChan := make(chan string, 50)
	go func() {
		// 使用有超时时间以及上下文管理的命令执行 方便处理异常以及适配暂停任务
		utils.Tools.ExecuteCommandToChanWithTimeout(findomainExecPath, []string{"-t", data, "-q"}, resultChan, time.Duration(executionTimeout)*time.Minute, op.Ctx)
	}()
	count := 0 // 初始化计数器
	for line := range resultChan {
		if strings.TrimSpace(line) != "" {
			count += 1
			//resultDns := utils.DNS.QueryOne(line)
			//tmp := utils.DNS.DNSdataToSubdomainResult(resultDns)
			tmp := types.SubdomainResult{
				Host:  line,
				Type:  "A",
				Value: []string{}, // 需要进一步实现获取 IP 的逻辑
				IP:    []string{},
				Tags:  []string{"findomain-scan"},
			}
			var addrs []net.IP
			addrs, err := net.LookupIP(line)
			if err != nil {
			} else {
				for _, addr := range addrs {
					tmp.IP = append(tmp.IP, addr.String()) // 使用 addr.String() 转换为字符串
				}
			}
			op.ResultFunc(tmp)
		}
	}
	end := time.Now()
	duration := end.Sub(start)
	op.Log(fmt.Sprintf("Findomain scan for %s completed, found %d subdomains, time: %v", data, count, duration))
	return nil, nil
}
