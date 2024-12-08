// ehole-------------------------------------
// @file      : plugin.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/11/21 21:55
// -------------------------------------------

package plugin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"net"
	"net/url"
	"strconv"
	"strings"
)

func GetName() string {
	return "fofa"
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

var (
	APIKEY = "xxxxxxxx"
)

type Response struct {
	Error           bool       `json:"error"`
	ConsumedFPoint  int        `json:"consumed_fpoint"`
	RequiredFPoints int        `json:"required_fpoints"`
	Size            int        `json:"size"`
	Page            int        `json:"page"`
	Mode            string     `json:"mode"`
	Query           string     `json:"query"`
	Results         [][]string `json:"results"` // results 是一个包含数组的二维数组
}

// Execute 目标处理
// 带*的目标 "*.example.com"
// 不带*的目标
func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	data, ok := input.(string)
	if !ok {
		return nil, nil
	}
	parameter := op.Parameter
	size := 10
	if parameter != "" {
		args, err := utils.Tools.ParseArgs(parameter, "size")
		if err != nil {
		} else {
			for key, value := range args {
				if value != "" {
					switch key {
					case "size":
						size, _ = strconv.Atoi(value)
					}
				}
			}
		}
	}
	if strings.Contains(data, ".") {
		data = strings.TrimPrefix(data, "http://")
		data = strings.TrimPrefix(data, "https://")
		var target string
		target = data
		// 域名或ip类
		// 查找 "*." 是否存在
		if strings.Contains(data, "*.") {
			// 查找 "*." 的位置
			startIndex := strings.Index(data, "*.") + len("*.")
			target = data[startIndex:]
		}
		target = "http://" + target
		parsedURL, err := url.Parse(target)
		if err != nil {
			op.Log(fmt.Sprintf("target %v parsedURL error: %v", data, err), "w")
			return nil, nil
		}
		tmpdata := strings.Split(parsedURL.Host, ":")
		target = tmpdata[0]
		ip := net.ParseIP(target)
		query := ""
		if ip != nil {
			query = fmt.Sprintf("ip=\"%s\"", target)
		} else {
			query = fmt.Sprintf("host=\"%s\"", target)
		}
		if query != "" {
			encoded := base64.StdEncoding.EncodeToString([]byte(query))
			urlRaw := fmt.Sprintf("https://fofa.info/api/v1/search/all?&key=%v&qbase64=%v&size=%v&fields=ip,port,host,protocol,icp,title", APIKEY, encoded, size)
			res, err := utils.Requests.HttpGetByte(urlRaw)
			if err != nil {
				op.Log(fmt.Sprintf("get target %v error: %v", query, err), "w")
				return nil, err
			}
			reader := bytes.NewReader(res)
			decoder := json.NewDecoder(reader)
			var result Response
			if err := decoder.Decode(&result); err != nil {
				op.Log(fmt.Sprintf("get target %v json decode error: %v", query, err), "w")
				return nil, err
			}
			if result.Error {
				op.Log(fmt.Sprintf("get target %v result error is true", query), "w")
				return nil, err
			}
			op.Log(fmt.Sprintf("target %v search：%v fofa All size: %v search result: %v required_fpoints:%v", data, query, result.Size, len(result.Results), result.RequiredFPoints))
			for _, r := range result.Results {
				// 假设每个 result 是一个包含 3 个字符串的数组
				// ip,port,domain,host,protocol,icp,title
				if len(r) == 6 {
					rawIp := r[0]
					port := r[1]
					host := r[2]
					protocol := r[3]
					icp := r[4]
					title := r[5]
					createURL := host
					if !strings.Contains(host, "http://") && !strings.Contains(host, "https://") {
						createURL = "http://" + host
					}
					assetDomain := host
					parsedURL, err = url.Parse(createURL)
					if err == nil {
						assetDomain = parsedURL.Host
					}
					// 构建asset资产
					if protocol == "http" || protocol == "https" {
						assetHttp := types.AssetHttp{
							Type:    "http",
							Port:    port,
							IP:      rawIp,
							Title:   title,
							URL:     host,
							Service: protocol,
							Tags:    []string{"FOFA"},
						}
						assetHttp.Host = assetDomain
						if icp != "" {
							assetHttp.Tags = append(assetHttp.Tags, fmt.Sprintf("icp:%v", icp))
						}
						op.ResultFunc(assetHttp)
					} else {
						assetOther := types.AssetOther{
							Type:    "other",
							Service: protocol,
							Port:    port,
							IP:      rawIp,
							Tags:    []string{"FOFA"},
						}
						assetOther.Host = assetDomain
						if icp != "" {
							assetOther.Tags = append(assetOther.Tags, fmt.Sprintf("icp:%v", icp))
						}
						op.ResultFunc(assetOther)
					}
					// 构建子域名
					subdomain := types.SubdomainResult{
						Host: assetDomain,
						IP:   []string{rawIp},
						Tags: []string{"FOFA"},
					}
					op.ResultFunc(subdomain)

					// 构建端口信息
					domainSkip := types.PortAlive{
						Host: assetDomain,
						IP:   rawIp,
						Port: port,
					}
					op.ResultFunc(domainSkip)
				}
			}
		}
	} else {
		// 非域名ip类，比如公司名称

	}

	return nil, nil
}
