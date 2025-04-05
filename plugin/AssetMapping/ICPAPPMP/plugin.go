// ICPQuery-------------------------------------
// @file      : plugin.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/4/4 17:51
// -------------------------------------------

package plugin

import (
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/global"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/results"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func GetName() string {
	return "ICPAPPMP"
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

var MIITAPI = ""

func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	rootDomainResult, ok := input.(types.RootDomain)
	if !ok {
		return nil, nil
	}
	parameter := op.Parameter
	// outputDirName := fmt.Sprintf("bbot_result_%s", targetMD5)

	if parameter != "" {
		args, err := utils.Tools.ParseArgs(parameter, "et", "finger", "common", "bak", "d")
		if err != nil {
			op.Log(fmt.Sprintf("parse args error: %v", err), "e")
		} else {
			for key, value := range args {
				if value != "" {
					switch key {

					default:
						continue
					}
				}

			}
		}
	}

	if rootDomainResult.ICP != "" {
		tmpIcp := LastSplitOnce(rootDomainResult.ICP, "-")
		if rootDomainResult.Company == "" {
			value, exists := global.TmpCustomMapParameter.Load(tmpIcp)
			if exists {
				strValue, _ := value.(string)
				rootDomainResult.Company = strValue
			}
		} else {
			global.TmpCustomMapParameter.Store(tmpIcp, rootDomainResult.Company)
		}
	}
	// 以下GetIcp接口是通过https://www.beianx.cn/进行查询
	if rootDomainResult.ICP == "" || rootDomainResult.Company == "" {
		// icp为空查询domain的icp
		// 去除重复查询的icp
		locakKey := "duplicates:" + op.TaskId + ":icpdomain:" + rootDomainResult.Domain
		keyRedis := "duplicates:" + op.TaskId + ":icpdomain"
		valueRedis := rootDomainResult.Domain
		flag := results.Duplicate.Custom(locakKey, keyRedis, valueRedis)
		if flag {
			// 不重复
			res := GetIcp(rootDomainResult.Domain)
			if len(res) != 0 {
				rootDomainResult.ICP = res[0].ICP
				rootDomainResult.Company = res[0].Company
				global.TmpCustomMapParameter.Store(LastSplitOnce(rootDomainResult.ICP, "-"), rootDomainResult.Company)
			} else {
				// 如果beian.cn没有找到信息，尝试使用miit接口查询
				if GetMiitApi() != "" {
					miitWebRes, err := GetWebByMiit(rootDomainResult.Domain, map[string]bool{})
					if err == nil {
						// 如果没有报错
						if len(miitWebRes) != 0 {
							rootDomainResult.ICP = miitWebRes[0].ICP
							rootDomainResult.Company = miitWebRes[0].Company
						}
					}
				}

			}
		} else {
			//域名在当前任务中已经查过了 后边也不用走了 直接返回
			return nil, nil
		}
	}
	var domainList map[string]bool
	domainList[rootDomainResult.Domain] = true
	// 将初始结果发送到结果处理
	op.ResultFunc(rootDomainResult)
	var allRootDomain []ICPinfo
	tmpIcp := LastSplitOnce(rootDomainResult.ICP, "-")
	// 根据icp查询更多的根域名
	if rootDomainResult.ICP != "" {
		// 如果icp不为空的话  根据icp 查更多的根域名
		locakKey := "duplicates:" + op.TaskId + ":icp:" + tmpIcp
		keyRedis := "duplicates:" + op.TaskId + ":icp"
		valueRedis := tmpIcp
		flag := results.Duplicate.Custom(locakKey, keyRedis, valueRedis)
		if flag {
			// 该icp没有查询过
			allRootDomain = GetIcp(tmpIcp)
			if len(allRootDomain) != 0 {
				// 根据icp查询根域名
				for _, r := range allRootDomain {
					if r.Domain != rootDomainResult.Domain {
						op.ResultFunc(types.RootDomain{Domain: r.Domain, ICP: r.ICP, Company: r.Company})
						domainList[r.Domain] = true
					}
				}
			}
		}
	}
	// 由于app、小程序接口依赖于miit接口 所以如果miit接口是空的就返回，据了解零零信安的接口也可以查询app和小程序 可以通过rootDomainResult的公司名进行查找
	// 我这没有vip，期待其他师傅补充零零信安的接口
	if GetMiitApi() == "" {
		return nil, nil
	}
	// 根据icp查找根域名、app、小程序 通过miit接口查询
	if rootDomainResult.ICP != "" {
		locakKey := "duplicates:" + op.TaskId + ":miit:" + tmpIcp
		keyRedis := "duplicates:" + op.TaskId + ":miit"
		valueRedis := tmpIcp
		flag := results.Duplicate.Custom(locakKey, keyRedis, valueRedis)
		if flag {
			// 该公司没有通过miit查询过, 根据icp查询更多的根域名
			miitWebRes, err := GetWebByMiit(tmpIcp, domainList)
			if err != nil {
				return nil, err
			}
			for _, r := range miitWebRes {
				op.ResultFunc(types.RootDomain{Domain: r.Domain, ICP: r.ICP, Company: r.Company})
			}
			// 查询app
			miitAppRes, err := GetAppByMiit(tmpIcp)
			if err != nil {
				return nil, err
			}
			for _, r := range miitAppRes {
				op.ResultFunc(r)
			}
			// 查询小程序
			miitMappRes, err := MiitMappApi(tmpIcp)
			if err != nil {
				return nil, err
			}
			for _, r := range miitMappRes {
				op.ResultFunc(r)
			}
		}
	}

	return nil, nil
}

func GetMiitApi() string {
	return MIITAPI
}

type MiitWebResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Params  Params `json:"params"`
}

type Params struct {
	EndRow           int        `json:"endRow"`
	FirstPage        int        `json:"firstPage"`
	HasNextPage      bool       `json:"hasNextPage"`
	HasPreviousPage  bool       `json:"hasPreviousPage"`
	IsFirstPage      bool       `json:"isFirstPage"`
	IsLastPage       bool       `json:"isLastPage"`
	LastPage         int        `json:"lastPage"`
	List             []UnitInfo `json:"list"`
	NavigatePages    int        `json:"navigatePages"`
	NavigatePageNums []int      `json:"navigatepageNums"`
	NextPage         int        `json:"nextPage"`
	PageNum          int        `json:"pageNum"`
	PageSize         int        `json:"pageSize"`
	Pages            int        `json:"pages"`
	PrePage          int        `json:"prePage"`
	Size             int        `json:"size"`
	StartRow         int        `json:"startRow"`
	Total            int        `json:"total"`
}

type UnitInfo struct {
	CityID           int    `json:"cityId"`
	CountyID         int    `json:"countyId"`
	DataID           int    `json:"dataId"`
	LeaderName       string `json:"leaderName"`
	MainID           string `json:"mainId"`
	MainLicence      string `json:"mainLicence"`
	MainUnitAddress  string `json:"mainUnitAddress"`
	MainUnitCertNo   string `json:"mainUnitCertNo"`
	MainUnitCertType int    `json:"mainUnitCertType"`
	NatureID         int    `json:"natureId"`
	NatureName       string `json:"natureName"`
	ProvinceID       int    `json:"provinceId"`
	ServiceID        int    `json:"serviceId"`
	ServiceLicence   string `json:"serviceLicence"`
	ServiceName      string `json:"serviceName"`
	ServiceType      int    `json:"serviceType"`
	UnitName         string `json:"unitName"`
	UpdateRecordTime string `json:"updateRecordTime"`
	Version          string `json:"version"`
	Domain           string `json:"domain"`
}

func MiitApi(str string, pageSize int, retryCount int, tp string) (MiitWebResponse, error) {
	// 判断是否已用完重试次数
	if retryCount <= 0 {
		return MiitWebResponse{}, fmt.Errorf("max retries reached")
	}

	api := GetMiitApi() + fmt.Sprintf("/query/%v?search=", tp) + str + "&pageSize=40" + fmt.Sprintf("&pageNum=%v", pageSize)
	timeout := 5 * time.Second
	maxRetries := 3
	retryInterval := 2 * time.Second
	headers := map[string]string{}
	proxy := ""
	// 获取第一页数据 获取总页数
	resp, err := utils.Requests.HttpGetWithRetry(api,
		timeout,       // 超时时间
		maxRetries,    // 最大重试次数
		retryInterval, // 重试间隔
		headers,       // 请求头
		proxy,         // 代理地址
	)
	var miitWebResponse MiitWebResponse
	if resp.StatusCode != 200 {
		logger.SlogWarnLocal(fmt.Sprintf("[%v] GetWebByMiit %v error page 1, resp.StatusCode %v", GetName(), str, resp.StatusCode))
		return miitWebResponse, fmt.Errorf("not 200")
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.SlogWarnLocal(fmt.Sprintf("[%v] GetWebByMiit %v Failed to read body %v", GetName(), str, err))
		return miitWebResponse, err
	}
	err = json.Unmarshal(body, &miitWebResponse)
	if err != nil {
		logger.SlogWarnLocal(fmt.Sprintf("[%v] GetWebByMiit %v Error unmarshalling response: %v", GetName(), str, err))
		return miitWebResponse, err
	}

	// 如果code不是200，进行递归调用并减少retryCount
	if miitWebResponse.Code != 200 {
		logger.SlogWarnLocal(fmt.Sprintf("[%v] GetWebByMiit %v code %v, retrying %v", GetName(), str, miitWebResponse.Code, retryCount-1))
		return MiitApi(str, pageSize, retryCount-1, tp)
	}

	return miitWebResponse, nil
}

func GetAppByMiit(str string) ([]types.APP, error) {
	var result []types.APP
	// 获取第一页数据 获取总页数
	res, err := MiitApi(str, 1, 5, "app")
	if err != nil {
		logger.SlogWarnLocal(fmt.Sprintf("[%v] GetWebByMiit %v error page 1: %v", GetName(), str, err))
		return result, err
	}
	if len(res.Params.List) != 0 {
		for _, r := range res.Params.List {
			result = append(result, types.APP{ICP: r.ServiceLicence, Name: r.ServiceName, Company: r.UnitName})
		}
	}
	for i := 2; i <= res.Params.Pages; i++ {
		res, err = MiitApi(str, i, 5, "app")
		if err != nil {
			logger.SlogWarnLocal(fmt.Sprintf("[%v] GetWebByMiit %v error page 1: %v", GetName(), str, err))
			return result, err
		}
		if len(res.Params.List) != 0 {
			for _, r := range res.Params.List {
				result = append(result, types.APP{ICP: r.ServiceLicence, Name: r.ServiceName, Company: r.UnitName})
			}
		}
	}
	return result, nil
}

func MiitMappApi(str string) ([]types.MP, error) {
	var result []types.MP
	// 获取第一页数据 获取总页数
	res, err := MiitApi(str, 1, 5, "mapp")
	if err != nil {
		logger.SlogWarnLocal(fmt.Sprintf("[%v] GetWebByMiit %v error page 1: %v", GetName(), str, err))
		return result, err
	}
	if len(res.Params.List) != 0 {
		for _, r := range res.Params.List {
			result = append(result, types.MP{ICP: r.ServiceLicence, Name: r.ServiceName, Company: r.UnitName})
		}
	}
	for i := 2; i <= res.Params.Pages; i++ {
		res, err = MiitApi(str, i, 5, "app")
		if err != nil {
			logger.SlogWarnLocal(fmt.Sprintf("[%v] GetWebByMiit %v error page 1: %v", GetName(), str, err))
			return result, err
		}
		if len(res.Params.List) != 0 {
			for _, r := range res.Params.List {
				result = append(result, types.MP{ICP: r.ServiceLicence, Name: r.ServiceName, Company: r.UnitName})
			}
		}
	}
	return result, nil
}

func GetWebByMiit(str string, domainList map[string]bool) ([]ICPinfo, error) {
	var result []ICPinfo
	// 获取第一页数据 获取总页数
	res, err := MiitApi(str, 1, 5, "web")
	if err != nil {
		logger.SlogWarnLocal(fmt.Sprintf("[%v] GetWebByMiit %v error page 1: %v", GetName(), str, err))
		return result, err
	}
	if len(res.Params.List) != 0 {
		for _, r := range res.Params.List {
			_, e := domainList[r.Domain]
			if !e {
				result = append(result, ICPinfo{Domain: r.Domain, ICP: r.ServiceLicence, Company: r.UnitName})
				domainList[r.Domain] = true
			}
		}
	}
	for i := 2; i <= res.Params.Pages; i++ {
		res, err = MiitApi(str, i, 5, "web")
		if err != nil {
			logger.SlogWarnLocal(fmt.Sprintf("[%v] GetWebByMiit %v error page 1: %v", GetName(), str, err))
			return result, err
		}
		if len(res.Params.List) != 0 {
			for _, r := range res.Params.List {
				_, e := domainList[r.Domain]
				if !e {
					result = append(result, ICPinfo{Domain: r.Domain, ICP: r.ServiceLicence, Company: r.UnitName})
					domainList[r.Domain] = true
				}
			}
		}
	}
	return result, nil
}

func LastSplitOnce(s, sep string) (before string) {
	idx := strings.LastIndex(s, sep)
	if idx == -1 {
		return s
	}
	return s[:idx]
}

type ICPinfo struct {
	Company string
	ICP     string
	Domain  string
}

var (
	trRe      = regexp.MustCompile(`<tr>([\s\S]+?)</tr>`)
	companyRe = regexp.MustCompile(`<a href="/company/[^"]+">([^<]+)</a>`)
	icpRe     = regexp.MustCompile(`<td class="align-middle" nowrap="nowrap">\s*([^<]+?)\s*</td>`)
	websiteRe = regexp.MustCompile(`<a href="\/seo\/([^"]+)">([^<]+)<\/a>\s*`)
)

func GetIcp(domain string) []ICPinfo {
	url := fmt.Sprintf("https://www.beianx.cn/search/%v", domain)
	timeout := 5 * time.Second
	maxRetries := 3
	retryInterval := 2 * time.Second
	// 尝试请求，最多重试 maxRetries 次
	var resp *http.Response
	var err error
	headers := map[string]string{}
	proxy := ""
	resp, err = utils.Requests.HttpGetWithRetry(
		url,
		timeout,       // 超时时间
		maxRetries,    // 最大重试次数
		retryInterval, // 重试间隔
		headers,       // 请求头
		proxy,         // 代理地址
	)
	var result []ICPinfo
	// 检查是否成功发起请求
	if err != nil {
		logger.SlogWarnLocal(fmt.Sprintf("[%v] domain %v Failed to make request after %d attempts: %v", GetName(), domain, maxRetries, err))
		return result
	}
	if resp.StatusCode != 200 {
		logger.SlogWarnLocal(fmt.Sprintf("[%v] domain %v Failed to query %v", GetName(), domain, resp.Status))
		return result
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.SlogWarnLocal(fmt.Sprintf("[%v] domain %v Failed to read body %v", GetName(), domain, err))
		return result
	}
	html := string(body)
	trMatches := trRe.FindAllStringSubmatch(html, -1)
	for _, tr := range trMatches {
		trContent := tr[1]
		if strings.Contains(trContent, "主办单位名称") {
			continue
		}
		// 提取网站地址
		website := ""
		if websiteMatch := websiteRe.FindStringSubmatch(trContent); len(websiteMatch) > 2 {
			website = strings.TrimSpace(websiteMatch[2])
		}
		if website == "" {
			continue
		}
		// 提取公司名称
		company := ""
		if companyMatch := companyRe.FindStringSubmatch(trContent); len(companyMatch) > 1 {
			company = strings.TrimSpace(companyMatch[1])
		}

		// 提取ICP备案号
		icp := ""
		if icpMatch := icpRe.FindStringSubmatch(trContent); len(icpMatch) > 1 {
			icp = strings.TrimSpace(icpMatch[1])
		}
		website = strings.TrimPrefix("www.", website)
		result = append(result, ICPinfo{
			Domain:  website,
			Company: company,
			ICP:     icp,
		})
	}
	return result
}
