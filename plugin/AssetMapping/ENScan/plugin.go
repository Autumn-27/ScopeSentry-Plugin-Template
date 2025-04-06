// ENScan-------------------------------------
// @file      : plugin.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/4/6 16:08
// -------------------------------------------

package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/global"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/notification"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetName() string {
	return "ENScan"
}

var ConfigStr = `version: 0.5
app:
  miit_api: ''          # HG-ha的ICP_Query (非狼组维护，团队成员请使用内部版本)
cookies:
  aiqicha: ''           # 爱企查   Cookie
  tianyancha: ''        # 天眼查   Cookie
  tycid: ''        		# 天眼查   CApi ID(capi.tianyancha.com)
  auth_token: ''        # 天眼查   Token (capi.tianyancha.com)
  qcc: ''               # 企查查   Cookie
  qcctid: '' 			# 企查查   TID console.log(window.tid)
  aldzs: ''             # 阿拉丁   Cookie
  xlb: ''               # 小蓝本   Token
  qimai: ''             # 七麦数据 Cookie
`

func Install() error {
	// 安装 enscan 到指定目录下面
	enscanPath := filepath.Join(global.ExtDir, "enscan")
	if err := os.MkdirAll(enscanPath, os.ModePerm); err != nil {
		logger.SlogError(fmt.Sprintf("Failed to create ksubdomain folder:", err))
		return err
	}
	osType := runtime.GOOS
	var enscanURL string
	var fileName string
	var zipName string
	switch osType {
	case "windows":
		enscanURL = "https://github.com/wgpsec/ENScan_GO/releases/download/v1.2.1/enscan-v1.2.1-windows-amd64.zip"
		fileName = "enscan.exe"
		zipName = "enscan.zip"
	case "linux":
		enscanURL = "https://github.com/wgpsec/ENScan_GO/releases/download/v1.2.1/enscan-v1.2.1-linux-amd64.tar.gz"
		fileName = "enscan"
		zipName = "enscan.tar.gz"
	}
	enscanExecPath := filepath.Join(enscanPath, fileName)
	if _, err := os.Stat(enscanExecPath); os.IsNotExist(err) {
		downloadPath := filepath.Join(global.ExtDir, "enscan", zipName) // 临时下载路径
		success, err := utils.Tools.HttpGetDownloadFile(enscanURL, downloadPath)
		if err != nil || !success {
			logger.SlogErrorLocal(fmt.Sprintf("Failed to download enscan: %v", err))
			return err
		}

		logger.SlogInfo("enscan Download successful")

		err = utils.Tools.UnzipFile(downloadPath, enscanPath)
		if err != nil {
			fmt.Printf("Failed to extract enscan: %v", err)
			return nil
		}
		switch osType {
		case "windows":
			os.Rename(filepath.Join(enscanPath, "enscan-v1.2.1-windows-amd64.exe"), enscanExecPath)
		case "linux":
			os.Rename(filepath.Join(enscanPath, "enscan-v1.2.1-linux-amd64"), enscanExecPath)
			enscanExecPath := filepath.Join(enscanPath, "enscan")
			err = os.Chmod(enscanExecPath, 0755)
			if err != nil {
				fmt.Sprintf("Failed to set permissions: %v", err)
				return nil
			}
		}
		defer utils.Tools.DeleteFile(downloadPath)
		err = utils.Tools.EnsureDir(filepath.Join(enscanPath, "result"))
		if err != nil {
			logger.SlogErrorLocal(fmt.Sprintf("Failed to EnsureDir: %v", err))
			return err
		}
	}
	utils.Tools.WriteContentFile(filepath.Join(enscanPath, "config.yaml"), ConfigStr)
	logger.SlogInfo("enscan installed successfully")
	return nil
}

func Check() error {
	return nil
}

func Uninstall() error {
	return nil
}

func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	companyTarget, ok := input.(types.Company)
	if !ok {
		return nil, nil
	}
	if companyTarget.Name == "" {
		return nil, nil
	}
	parameter := op.Parameter
	var (
		tp, field, delay, invest string
	)
	invest = ""
	delay = ""
	field = ""
	tp = ""
	executionTimeout := 10
	arg := []string{}
	if parameter != "" {
		args, err := utils.Tools.ParseArgs(parameter, "type", "field", "delay", "invest", "et")
		if err != nil {
			op.Log(fmt.Sprintf("parse args error: %v", err), "e")
		} else {
			for key, value := range args {
				if value != "" {
					switch key {
					case "type":
						tp = value
					case "field":
						field = value
					case "delay":
						delay = value
					case "invest":
						invest = value
					case "et":
						executionTimeout, _ = strconv.Atoi(value)
					default:
						continue
					}
				}
			}
		}
	}
	arg = append(arg, "-n")
	arg = append(arg, companyTarget.Name)
	if tp != "" {
		arg = append(arg, "-type")
		arg = append(arg, tp)
	}
	if field != "" {
		arg = append(arg, "-field")
		arg = append(arg, field)
	}
	if delay != "" {
		arg = append(arg, "-delay")
		arg = append(arg, delay)
	}
	if invest != "" {
		arg = append(arg, "-invest")
		arg = append(arg, invest)
	}
	arg = append(arg, "-json")
	resultTmpPath := utils.Tools.GenerateRandomString(6)
	resultPath := filepath.Join(global.ExtDir, "enscan", "result", resultTmpPath)
	utils.Tools.EnsureDir(resultPath)
	arg = append(arg, "--out-dir")
	arg = append(arg, resultPath)
	osType := runtime.GOOS
	var fileName string
	switch osType {
	case "windows":
		fileName = "enscan.exe"
	case "linux":
		fileName = "enscan"
	}
	enscanExecPath := filepath.Join(global.ExtDir, "enscan", fileName)
	resultChan := make(chan string, 50)
	// 使用 WithCancel 创建一个新的上下文
	newCtx, cancel := context.WithCancel(op.Ctx)
	defer cancel()
	go func() {
		utils.Tools.ExecuteCommandToChanWithTimeout(enscanExecPath, arg, resultChan, time.Duration(executionTimeout)*time.Minute, newCtx)
	}()
	flag := 0
	for line := range resultChan {
		if strings.TrimSpace(line) != "" {
			logger.SlogInfoLocal(fmt.Sprintf("[Plguin %v] %v", GetName(), line))
			if strings.Contains(line, "10秒后重试") {
				logger.SlogInfoLocal(fmt.Sprintf("[Plguin %v] 需要重试", GetName()))
				tmp := line
				notification.FlushBuffer("Warn", &tmp)
				flag += 1
				if flag > 10 {
					logger.SlogWarn(fmt.Sprintf("[Plguin %v] 超时10次结束执行 %v", GetName(), companyTarget.Name))
					cancel()
					return nil, nil
				}
			}
		}
	}
	path, err := GetFirstJSONFilePath(resultPath)
	if err != nil {
		return nil, err
	}
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		logger.SlogWarnLocal(fmt.Sprintf("[Plugin %v]Error opening file: %v", GetName(), err))
		return nil, err
	}
	defer file.Close()

	// 读取文件内容
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		logger.SlogWarnLocal(fmt.Sprintf("[Plugin %v]Error reading file: %v", GetName(), err))
		return nil, err
	}

	// 定义结构体变量
	var data Root

	// 解析 JSON 数据
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		logger.SlogWarnLocal(fmt.Sprintf("[Plugin %v]Error unmarshalling JSON: %v", GetName(), err))
		return nil, err
	}
	company := ""
	if len(data.EnterpriseInfo) != 0 {
		company = data.EnterpriseInfo[0].Name
	}
	// rootDomain
	for _, icp := range data.ICP {
		tmpIcp := types.RootDomain{
			ICP:    icp.ICP,
			Domain: icp.Domain,
		}
		if icp.From != "" {
			tmpIcp.Company = icp.From
		} else {
			tmpIcp.Company = company
		}
		op.ResultFunc(tmpIcp)
	}
	// app
	for _, app := range data.App {
		tmpAPP := types.APP{Name: app.Name, BundleID: app.BundleID, Url: app.Link, Category: app.Category}
		if app.From != "" {
			tmpAPP.Company = app.From
		} else {
			tmpAPP.Company = company
		}
		op.ResultFunc(tmpAPP)
	}
	// 小程序
	for _, mp := range data.WxApp {
		tmpMP := types.MP{
			Name:     mp.Name,
			Category: mp.Category,
		}
		if mp.From != "" {
			tmpMP.Company = mp.From
		} else {
			tmpMP.Company = company
		}
		op.ResultFunc(tmpMP)
	}
	return nil, nil
}

// GetFirstJSONFilePath 从指定目录中返回第一个 .json 文件的完整路径
func GetFirstJSONFilePath(dirPath string) (string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("读取目录失败: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			fullPath := filepath.Join(dirPath, entry.Name())
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("目录中没有找到 json 文件")
}

type Root struct {
	App            []AppInfo        `json:"app"`             // 应用信息列表
	ICP            []ICPInfo        `json:"icp"`             // ICP备案信息列表
	WxApp          []WxAppInfo      `json:"wx_app"`          // 微信小程序信息列表
	EnterpriseInfo []EnterpriseInfo `json:"enterprise_info"` // 企业信息列表
}

// 应用信息结构体
type AppInfo struct {
	BundleID    string `json:"bundle_id"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Extra       string `json:"extra"`
	From        string `json:"from"`
	Link        string `json:"link"`
	Logo        string `json:"logo"`
	Market      string `json:"market"`
	Name        string `json:"name"`
	UpdateAt    string `json:"update_at"`
	Version     string `json:"version"`
}

// ICP备案信息结构体
type ICPInfo struct {
	CompanyName string `json:"company_name"`
	Domain      string `json:"domain"`
	Extra       string `json:"extra"`
	From        string `json:"from"`
	ICP         string `json:"icp"`
	Website     string `json:"website"`
	WebsiteName string `json:"website_name"`
}

// 微信小程序信息结构体
type WxAppInfo struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Logo     string `json:"logo"`
	QRCode   string `json:"qrcode"`
	ReadNum  string `json:"read_num"`
	From     string `json:"from"`
}

// 企业信息结构体
type EnterpriseInfo struct {
	Address           string `json:"address"`            // 公司地址
	Email             string `json:"email"`              // 邮箱
	Extra             string `json:"extra"`              // 额外信息
	From              string `json:"from"`               // 来源
	IncorporationDate string `json:"incorporation_date"` // 成立日期
	LegalPerson       string `json:"legal_person"`       // 法人
	Name              string `json:"name"`               // 企业名称
	Phone             string `json:"phone"`              // 联系电话
	PID               string `json:"pid"`                // 企业ID
	RegCode           string `json:"reg_code"`           // 注册号/统一社会信用代码
	RegisteredCapital string `json:"registered_capital"` // 注册资本
	Scope             string `json:"scope"`              // 经营范围
	Status            string `json:"status"`             // 企业状态
}
