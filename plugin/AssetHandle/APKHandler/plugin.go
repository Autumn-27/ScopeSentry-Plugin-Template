// APKHandler-------------------------------------
// @file      : plugin.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/4/9 20:21
// -------------------------------------------

package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/global"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func GetName() string {
	return "APKHandler"
}

var (
	// APKHandlerPath 是 APKHandler 的根目录
	APKHandlerPath = filepath.Join(global.ExtDir, "APKHandler")

	// AppPath 是 APKHandler/app，用于存放原始 app 文件
	AppPath = filepath.Join(APKHandlerPath, "app")

	// TmpPath 是 APKHandler/tmp，用于存放 app 解压后的内容（如 dex）
	TmpPath = filepath.Join(APKHandlerPath, "tmp")

	// ApktoolPath 是 APKHandler/apktool，用于存放 apktool 反编译的结果
	ApktoolPath = filepath.Join(APKHandlerPath, "apktool")

	// JavaPath 是 APKHandler/java，用于存放 dex 转 jar 后再反编译成的 Java 源码
	JavaPath = filepath.Join(APKHandlerPath, "java")
)

func Install() error {
	createDir := func(path string) error {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			logger.SlogError(fmt.Sprintf("[Plugin %v]Failed to create %s folder: %v", GetName(), path, err))
			return err
		}
		return nil
	}
	if err := createDir(APKHandlerPath); err != nil {
		return err
	}
	if err := createDir(AppPath); err != nil {
		return err
	}
	if err := createDir(TmpPath); err != nil {
		return err
	}
	if err := createDir(ApktoolPath); err != nil {
		return err
	}
	if err := createDir(JavaPath); err != nil {
		return err
	}

	// 检查java环境
	checkJavaEnvironment()

	return nil
}

func checkJavaEnvironment() bool {
	// 执行 java -version 命令
	cmd := exec.Command("java", "-version")
	stderr, err := cmd.CombinedOutput()

	// 如果执行失败，返回 false
	if err != nil {
		fmt.Println("Error executing java command:", err)
		return false
	}

	// 检查输出中是否包含 "java version"
	if strings.Contains(string(stderr), "java version") {
		return true
	}

	return false
}

func Check() error {
	return nil
}

func Uninstall() error {
	return nil
}

func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	appResult, ok := input.(*types.APP)
	if !ok {
		return nil, nil
	}
	downloadUrl := ""
	appName := ""
	var err error
	if appResult.BundleID != "" {
		// id不为空 直接获取下载链接 进行下载
		nameFlag := false
		if appResult.Name == "" {
			nameFlag = true
		}
		downloadUrl, appName, err = GetApkpureDownloadUrl(appResult.BundleID, nameFlag)
		if err != nil {
			return nil, err
		}
		if nameFlag {
			appResult.Name = appName
		}
	} else {
		// 如果id为空 根据app名称获取id 然后根据id获取下载链接
		if appResult.Name == "" {
			return nil, nil
		}
		id := GetIdByName(appResult.Name)
		appResult.BundleID = id
		downloadUrl, _, err = GetApkpureDownloadUrl(id, false)
		if err != nil {
			return nil, err
		}
	}
	// 获取到download url 下载apk
	if downloadUrl == "" {
		return nil, nil
	}
	return nil, nil
}

func GetIdByName(name string) string {
	id := HuaweiGetId(name)
	if id != "" {
		//logger.SlogInfoLocal(fmt.Sprintf("[Plugin %v]app %v HuaweiGetId : %v", GetName(), name, id))
		return id
	} else {
		getId, err := TencentGetId(name)
		if err != nil {
			return ""
		}
		return getId
	}
}

func HuaweiGetId(name string) string {
	interfaceCode := ""
	var err error
	for i := 1; i < 5; i++ {
		interfaceCode, err = HuaweiGetInterfaceCode()
		if err != nil {
			//logger.SlogWarnLocal(fmt.Sprintf("[Plugin %v]app %v HuwaGetInterfaceCode error: %v", GetName(), name, err))
			time.Sleep(5 * time.Second)
			continue
		}
		if interfaceCode == "error" {
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	if interfaceCode == "" {
		return ""
	}
	id, err := HuaweiSearch(name, interfaceCode)
	if err != nil {
		//logger.SlogWarnLocal(fmt.Sprintf("[Plugin %v]app %v HuaweiSearch error: %v", GetName(), name, err))
	}
	return id
}

type HuaweiSearchResponse struct {
	StatKey    string   `json:"statKey"`
	TitleType  string   `json:"titleType"`
	LayoutData []Layout `json:"layoutData"`
}

type Layout struct {
	LayoutId   int       `json:"layoutId"`
	LayoutName string    `json:"layoutName"`
	DataList   []AppItem `json:"dataList"`
}

type AppItem struct {
	Name    string `json:"name"`
	Memo    string `json:"memo"`
	Icon    string `json:"icon"`
	AppId   string `json:"appid"`
	Version string `json:"appVersionName"`
	Package string `json:"package"`
}

func HuaweiSearch(name string, code string) (string, error) {
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://web-dra.hispace.dbankcloud.com/edge/uowap/index?method=internal.getTabDetail&serviceType=20&reqPageNum=1&uri=searchApp%%7C%v&maxResults=25&version=10.0.0&zone=&locale=zh", url.QueryEscape(name)), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Host", "web-dra.hispace.dbankcloud.com")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="127", "Not)A;Brand";v="99"`)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Interface-Code", code)
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.89 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Origin", "https://appgallery.huawei.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://appgallery.huawei.com/")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("resp not 200 is %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// 解析 JSON 到结构体
	var result HuaweiSearchResponse
	err = json.Unmarshal(bodyText, &result)
	if err != nil {
		return "", err
	}
	for _, lay := range result.LayoutData {
		for _, dat := range lay.DataList {
			if dat.Name == name {
				return dat.Package, nil
			}
		}
	}

	return "", nil
}

func HuaweiGetInterfaceCode() (string, error) {
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	var data = strings.NewReader(`{"params":{},"zone":"","locale":"zh"}`)
	req, err := http.NewRequest("POST", "https://web-drcn.hispace.dbankcloud.com/edge/webedge/getInterfaceCode", data)
	if err != nil {
	}
	req.Header.Set("Host", "web-drcn.hispace.dbankcloud.com")
	req.Header.Set("Content-Length", "37")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="127", "Not)A;Brand";v="99"`)
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.89 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Origin", "https://appgallery.huawei.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://appgallery.huawei.com/")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "error", fmt.Errorf("resp not 200 is %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("%s\n", bodyText)
	res := strings.Trim(string(bodyText), "\"")
	if strings.HasPrefix(res, "eyJh") {
		timeValue := ""
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "HWWAFSESTIME" {
				timeValue = cookie.Value
				break
			}
		}
		return fmt.Sprintf("%v_%v", res, timeValue), nil
	} else {
		return "error", nil
	}
}

type SearchListItem struct {
	PkgName string `json:"pkg_name"`
	AppID   string `json:"app_id"`
	Name    string `json:"name"`
}

type SearchListData struct {
	Name                string           `json:"name"`
	ShowCloudGameButton bool             `json:"showCloudGameButton"`
	ItemData            []SearchListItem `json:"itemData"`
	CardID              string           `json:"cardid"`
}

type DynamicCardResponse struct {
	Ret            int    `json:"ret"`
	RequestID      string `json:"requestId"`
	GalileoTraceID string `json:"galileoTraceId"`
	Data           struct {
		Components []struct {
			Data SearchListData `json:"data"`
		} `json:"components"`
	} `json:"data"`
}

type PageProps struct {
	CanonicalURL        string                 `json:"canonicalUrl"`
	Context             map[string]interface{} `json:"context"`
	DynamicCardResponse DynamicCardResponse    `json:"dynamicCardResponse"`
}

type Props struct {
	PageProps PageProps `json:"pageProps"`
}

var TencentDataRE = regexp.MustCompile(`<script[^>]*id="__NEXT_DATA__"[^>]*>([\s\S]*?)<\/script>`)

func TencentGetId(name string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://sj.qq.com/search?q=%E7%99%BE%E5%BA%A6", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	//req.Header.Set("cookie", "qq_domain_video_guid_verify=6dbd2b7f9a6ed81c; pgv_pvid=7146943478; _qimei_q36=; _qimei_h38=c71d16fd92c35ebd6f50d51c02000007118308; tvfe_boss_uuid=8960c33d4ae454b5; RK=GaFxUoPj9u; ptcz=27a7832087f5b717eacf0f28de9c919cbf4f9a7df206105b6f16ac18b14ee4a2; _qimei_fingerprint=b06f6d3c7aeff53afe89635dc2b9eaba; _clck=3900699078|1|fuo|0; YYB_HOME_UUID=80b10196-0942-499f-88bb-19885e18e34a; tgw_l7_route=1913e3d3b7747a3906262ef4833d1290; is_gray=0; Hm_lvt_bee22ad562886a0c3c9e70e97e67022f=1743761693,1743844128,1744042478,1744209403; HMACCOUNT=8ECD5B2769A472F3; Hm_lpvt_bee22ad562886a0c3c9e70e97e67022f=1744209408")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	matches := TencentDataRE.FindStringSubmatch(string(bodyText))
	if len(matches) < 2 {
		return "", nil
	}
	var propsWrapper struct {
		Props Props `json:"props"`
	}

	// 将 JSON 字符串解码到结构体
	err = json.Unmarshal([]byte(matches[1]), &propsWrapper)
	if err != nil {
		return "", nil
	}
	if propsWrapper.Props.PageProps.DynamicCardResponse.Data.Components != nil {
		for _, dy := range propsWrapper.Props.PageProps.DynamicCardResponse.Data.Components {
			if dy.Data.ItemData != nil {
				for _, it := range dy.Data.ItemData {
					if it.Name == name {
						return it.PkgName, nil
					}
				}
			}
		}
	}
	return "", nil
}

var apkPureDownloadURLRegex = regexp.MustCompile(`(X?APKJ)..(https?://(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*))`)

var apkNameRegex = regexp.MustCompile(`(?s)version_list........([\P{C}]+)`)

func GetApkpureDownloadUrl(id string, getName bool) (string, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.pureapk.com/m/v3/cms/app_version?hl=en-US&package_name=%v", id), nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "")
	req.Header.Set("x-cv", "3172501")
	req.Header.Set("x-sv", "29")
	req.Header.Set("x-abis", "arm64-v8a,armeabi-v7a,armeabi")
	req.Header.Set("x-gp", "1")
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	// 查找匹配的结果
	matches := apkPureDownloadURLRegex.FindStringSubmatch(string(bodyText))

	if len(matches) >= 3 {
		name := ""
		if getName {
			nameMatches := apkNameRegex.FindStringSubmatch(string(bodyText))
			if len(nameMatches) >= 2 {
				name = nameMatches[1]
			}
		}
		return matches[2], name, nil
	} else {
		return "", "", fmt.Errorf("Not Fund download url")
	}
}

func main() {
	//res, _ := TencentGetId("百度极速版")
	fmt.Println(GetApkpureDownloadUrl("com.instagram.android", true))
}

// 解压APK文件
func unzipAPK(apkPath, outputDir string) error {
	r, err := zip.OpenReader(apkPath)
	if err != nil {
		return fmt.Errorf("failed to open apk file: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(outputDir, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", fpath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory for file %s: %v", fpath, err)
		}

		destFile, err := os.Create(fpath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			destFile.Close()
			return fmt.Errorf("failed to open file %s: %v", f.Name, err)
		}

		_, err = io.Copy(destFile, rc)
		rc.Close()
		destFile.Close()
		if err != nil {
			return fmt.Errorf("failed to copy content of %s: %v", f.Name, err)
		}
	}
	return nil
}

// 使用 d2j-dex2jar 将 .dex 文件转换为 .jar
func dexToJar(dexFile, outputDir string) (string, error) {
	// 构建 .bat 路径（相对路径转绝对路径）
	batRelativePath := filepath.Join("tool", "dex-tools-v2.4", "d2j-dex2jar.bat")
	batAbsPath, err := filepath.Abs(batRelativePath)
	if err != nil {
		return "", fmt.Errorf("无法获取 .bat 的绝对路径: %v", err)
	}

	// 获取 .bat 所在目录，作为工作目录
	workDir := filepath.Dir(batAbsPath)

	// 拼接输出的 jar 文件路径
	jarOutputPath := filepath.Join(outputDir, filepath.Base(dexFile)+".jar")

	// 构建执行命令
	cmd := exec.Command(batAbsPath, dexFile, "-o", jarOutputPath, "--force")
	output, err := cmd.CombinedOutput()

	// 输出调试信息
	fmt.Println("执行 .bat 路径：", batAbsPath)
	fmt.Println("工作目录：", workDir)
	fmt.Println("输出 jar 路径：", jarOutputPath)

	if err != nil {
		return "", fmt.Errorf("d2j-dex2jar 执行失败: %v\n输出内容: %s", err, string(output))
	}

	fmt.Println("转换成功")
	return jarOutputPath, nil
}

// jarToJava 将 .jar 反编译为 .java 文件，输出到 result/java/xxx/
func jarToJava(jarFile string, javaBaseDir string) error {
	jarName := strings.TrimSuffix(filepath.Base(jarFile), ".jar")
	outputDir := filepath.Join(javaBaseDir, jarName)

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output dir %s: %v", outputDir, err)
	}

	cmd := exec.Command("java", "-jar", "tool/cfr-0.152.jar", jarFile, "--outputdir", outputDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("CFR 反编译失败: %v\n输出内容: %s", err, string(output))
	}
	return nil
}

// apkToolDecompile 使用 apktool 反编译 apk
func apkToolDecompile(apkPath, outputDir string) error {
	cmd := exec.Command("java", "-jar", "tool/apktool_2.11.1.jar", "d", apkPath, "-o", outputDir, "-f")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("apktool 反编译失败: %v\n输出内容: %s", err, string(output))
	}
	return nil
}

// processDexFiles 处理 dex 文件
func processDexFiles(tmpDir, javaOutputDir string) error {
	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return fmt.Errorf("failed to read tmpDir: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".dex") {
			dexPath := filepath.Join(tmpDir, file.Name())
			jarPath, err := dexToJar(dexPath, tmpDir)
			if err != nil {
				return fmt.Errorf("failed to convert %s to jar: %v", file.Name(), err)
			}

			err = jarToJava(jarPath, javaOutputDir)
			if err != nil {
				return fmt.Errorf("failed to convert jar to java for %s: %v", file.Name(), err)
			}
		}
	}
	return nil
}
