// APKHandler-------------------------------------
// @file      : plugin.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/4/9 20:21
// -------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func GetName() string {
	return "APKHandler"
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

func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	//appResult, ok := input.(*types.APP)
	//if !ok {
	//	return nil, nil
	//}

	return nil, nil
}

func GetIdByName(name string) string {
	id := HuaweiGetId(name)
	if id != "" {
		//logger.SlogInfoLocal(fmt.Sprintf("[Plugin %v]app %v HuaweiGetId : %v", GetName(), name, id))
		return name
	}
	return ""
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

func TencentGetId(name string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://sj.qq.com/search?q=%E7%99%BE%E5%BA%A6", nil)
	if err != nil {
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
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
	}
	matches := TencentDataRE.FindStringSubmatch(string(bodyText))
	if len(matches) < 2 {
		return ""
	}
	var propsWrapper struct {
		Props Props `json:"props"`
	}

	// 将 JSON 字符串解码到结构体
	err = json.Unmarshal([]byte(matches[1]), &propsWrapper)
	if err != nil {
		return ""
	}
	if propsWrapper.Props.PageProps.DynamicCardResponse.Data.Components != nil {
		for _, dy := range propsWrapper.Props.PageProps.DynamicCardResponse.Data.Components {
			if dy.Data.ItemData != nil {
				for _, it := range dy.Data.ItemData {
					if it.Name == name {
						return it.PkgName
					}
				}
			}
		}
	}
	return ""
}

func main() {
	res := TencentGetId("百度极速版")
	fmt.Println(res)
}
