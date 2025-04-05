// results-------------------------------------
// @file      : handler.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/9/19 19:10
// -------------------------------------------

package results

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/global"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/mongodb"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/notification"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type handler struct {
}

var Handler *handler

func InitializeHandler() {
	Handler = &handler{}
}

func (h *handler) GetAssetProject(host string) string {
	for _, p := range global.Projects {
		for _, t := range p.Target {
			if host == t {
				// 判断是否为黑名单
				for _, igTarget := range p.IgnoreList {
					if host == igTarget {
						return ""
					}
				}
				// 正则判断是否为黑名单
				for _, igRegex := range p.IgnoreRegexList {
					if igRegex.MatchString(host) {
						return ""
					}
				}
				return p.ID
			}
		}
	}
	return ""
}

func (h *handler) Subdomain(result *types.SubdomainResult) {
	var interfaceSlice interface{}
	rootDomain, err := utils.Tools.GetRootDomain(result.Host)
	if err != nil {
		logger.SlogInfoLocal(fmt.Sprintf("%v GetRootDomain error: %v", result.Host, err))
	}
	result.RootDomain = rootDomain
	if result.Time == "" {
		result.Time = utils.Tools.GetTimeNow()
	}
	result.Project = h.GetAssetProject(rootDomain)
	interfaceSlice = &result
	if global.NotificationConfig.SubdomainScan {
		NotificationMsg := fmt.Sprintf("%v - %v\n", result.Host, result.IP)
		notification.NotificationQueues["SubdomainScan"].Queue <- NotificationMsg
	}
	ResultQueues["SubdomainScan"].Queue <- interfaceSlice
}

func (h *handler) SubdomainTakeover(result *types.SubTakeResult) {
	var interfaceSlice interface{}
	rootDomain, err := utils.Tools.GetRootDomain(result.Input)
	if err != nil {
		logger.SlogInfoLocal(fmt.Sprintf("%v GetRootDomain error: %v", result.Input, err))
	}
	result.RootDomain = rootDomain
	result.Project = h.GetAssetProject(rootDomain)
	interfaceSlice = &result
	if global.NotificationConfig.SubdomainTakeoverNotification {
		NotificationMsg := fmt.Sprintf("Subdomain Takeover:\n%v - %v\n", result.Input, result.Cname)
		notification.NotificationQueues["SubdomainSecurity"].Queue <- NotificationMsg
	}
	ResultQueues["SubdomainSecurity"].Queue <- interfaceSlice
}

func (h *handler) AssetChangeLog(result *types.AssetChangeLog) {
	var interfaceSlice interface{}
	interfaceSlice = &result
	ResultQueues["AssetChangeLog"].Queue <- interfaceSlice
}

func (h *handler) AssetUpdate(id string, updateData interface{}) {
	// 资产比较少，并且需要实时获取历史记录，所以不采用多条插入的方式
	objectID, err := primitive.ObjectIDFromHex(id) // 将字符串 ID 转换为 ObjectID
	if err != nil {
		logger.SlogError(fmt.Sprintf("AssetUpdate %v ObjectIDFromHex error:%v", id, err))
		return
	}
	selector := bson.M{"_id": objectID} // 创建选择器，匹配指定的 ObjectID
	_, err = mongodb.MongodbClient.Update("asset", selector, bson.M{"$set": updateData})
	if err != nil {
		logger.SlogError(fmt.Sprintf("AssetUpdate %v error:%v", id, err))
	} // 使用 $set 更新字段
}

func (h *handler) AssetOtherInsert(result *types.AssetOther) {
	// 资产比较少，并且需要实时获取历史记录，所以不采用多条插入的方式
	var interfaceSlice interface{}
	rootDomain, err := utils.Tools.GetRootDomain(result.Host)
	if err != nil {
		logger.SlogInfoLocal(fmt.Sprintf("%v GetRootDomain error: %v", result.Host, err))
	}
	result.RootDomain = rootDomain
	if result.Time == "" {
		result.Time = utils.Tools.GetTimeNow()
	}
	result.Project = h.GetAssetProject(rootDomain)
	interfaceSlice = &result
	_, err = mongodb.MongodbClient.InsertOne("asset", interfaceSlice)
	if err != nil {
		logger.SlogError(fmt.Sprintf("AssetOtherInsert error:%v", err))
		return
	}
}

func (h *handler) AssetHttpInsert(result *types.AssetHttp) {
	// 资产比较少，并且需要实时获取历史记录，所以不采用多条插入的方式
	var interfaceSlice interface{}
	rootDomain, err := utils.Tools.GetRootDomain(result.Host)
	if err != nil {
		logger.SlogInfoLocal(fmt.Sprintf("%v GetRootDomain error: %v", result.Host, err))
	}
	result.RootDomain = rootDomain
	if result.Time == "" {
		result.Time = utils.Tools.GetTimeNow()
	}
	result.Project = h.GetAssetProject(rootDomain)
	interfaceSlice = &result
	_, err = mongodb.MongodbClient.InsertOne("asset", interfaceSlice)
	if err != nil {
		logger.SlogError(fmt.Sprintf("AssetOtherInsert error:%v", err))
		return
	}
}

func (h *handler) URL(result *types.UrlResult) {
	var interfaceSlice interface{}
	rootDomain, err := utils.Tools.GetRootDomain(result.Input)
	if err != nil {
		logger.SlogInfoLocal(fmt.Sprintf("%v GetRootDomain error: %v", result.Input, err))
	}
	result.RootDomain = rootDomain
	result.Project = h.GetAssetProject(rootDomain)
	if result.Time == "" {
		result.Time = utils.Tools.GetTimeNow()
	}
	// 创建result的副本，并将Body设置为空
	// 创建一个新的result对象，但Body设为空
	resultCopy := types.UrlResult{
		Input:      result.Input,
		Source:     result.Source,
		OutputType: result.OutputType,
		Output:     result.Output,
		Status:     result.Status,
		Length:     result.Length,
		Time:       result.Time,
		Project:    result.Project,
		TaskName:   result.TaskName,
		ResultId:   result.ResultId,
		RootDomain: result.RootDomain,
		Body:       "", // 设置Body为空
		Tags:       result.Tags,
		Ext:        result.Ext,
	}
	interfaceSlice = &resultCopy
	ResultQueues["URLScan"].Queue <- interfaceSlice
}

func (h *handler) Crawler(result *types.CrawlerResult) {
	var interfaceSlice interface{}
	rootDomain, err := utils.Tools.GetRootDomain(result.Url)
	if err != nil {
		logger.SlogInfoLocal(fmt.Sprintf("%v GetRootDomain error: %v", result.Url, err))
	}
	result.RootDomain = rootDomain
	if result.Time == "" {
		result.Time = utils.Tools.GetTimeNow()
	}
	result.Project = h.GetAssetProject(rootDomain)
	interfaceSlice = &result
	ResultQueues["WebCrawler"].Queue <- interfaceSlice
}

func (h *handler) Sensitive(result *types.SensitiveResult) {
	var interfaceSlice interface{}
	rootDomain, err := utils.Tools.GetRootDomain(result.Url)
	if err != nil {
		logger.SlogInfoLocal(fmt.Sprintf("%v GetRootDomain error: %v", result.Url, err))
	}
	result.RootDomain = rootDomain
	if result.Time == "" {
		result.Time = utils.Tools.GetTimeNow()
	}
	result.Project = h.GetAssetProject(rootDomain)
	interfaceSlice = &result
	if global.NotificationConfig.SensitiveNotification {
		NotificationMsg := fmt.Sprintf("Sensitive Scan:\n%v - %v\n", result.Url, result.SID)
		notification.NotificationQueues["URLSecurity"].Queue <- NotificationMsg
	}
	ResultQueues["SensitiveResult"].Queue <- interfaceSlice

}

func (h *handler) SensitiveBody(body string, md5 string) {
	selector := bson.M{"md5": md5}

	update := bson.M{
		"$set": bson.M{
			"body": body,
			"md5":  md5,
		},
	}
	op := types.BulkUpdateOperation{
		Selector: selector,
		Update:   update,
	}
	ResultQueues["SensitiveBody"].Queue <- op
}

func (h *handler) SensitiveUrl(url string, urlId string, bodyId string) {
	key := "SenU:" + urlId
	collectionName := "SensitiveUrl"
	selector := bson.M{"urlId": urlId}
	flag := Duplicate.DuplicateLocalCache(key)
	if flag {
		update := bson.M{
			"$set": bson.M{
				"url":    url,
				"urlId":  urlId,
				"bodyId": bodyId,
			},
		}
		// 调用 Upsert 方法，执行插入或更新操作
		_, err := mongodb.MongodbClient.Upsert(collectionName, selector, update)
		if err != nil {
			logger.SlogError(fmt.Sprintf("SensitiveUrl insert mongodb error:%v", err))
		}
	}
}

func (h *handler) Dir(result *types.DirResult) {
	var interfaceSlice interface{}
	rootDomain, err := utils.Tools.GetRootDomain(result.Url)
	if err != nil {
		logger.SlogInfoLocal(fmt.Sprintf("%v GetRootDomain error: %v", result.Url, err))
	}
	result.RootDomain = rootDomain
	result.Project = h.GetAssetProject(rootDomain)
	interfaceSlice = &result
	if global.NotificationConfig.DirScanNotification {
		NotificationMsg := ""
		if result.Msg != "" {
			NotificationMsg = fmt.Sprintf("%v - %v -%v\n", result.Url, result.Status, result.Msg)
		} else {
			NotificationMsg = fmt.Sprintf("%v - %v\n", result.Url, result.Status)
		}
		notification.NotificationQueues["DirScan"].Queue <- NotificationMsg
	}
	ResultQueues["DirScan"].Queue <- interfaceSlice
}

func (h *handler) Vulnerability(result *types.VulnResult) {
	var interfaceSlice interface{}
	rootDomain, err := utils.Tools.GetRootDomain(result.Url)
	if err != nil {
		logger.SlogInfoLocal(fmt.Sprintf("%v GetRootDomain error: %v", result.Url, err))
	}
	result.RootDomain = rootDomain
	if result.Time == "" {
		result.Time = utils.Tools.GetTimeNow()
	}
	result.Project = h.GetAssetProject(rootDomain)
	interfaceSlice = &result
	if global.NotificationConfig.VulNotification {
		NotificationMsg := ""
		if global.NotificationConfig.VulLevel != "" {
			if strings.Contains(strings.ToLower(global.NotificationConfig.VulLevel), strings.ToLower(result.Level)+",") {
				if result.Url != "" {
					NotificationMsg += fmt.Sprintf("%v-[%v]-[%v]-[%v]\n", result.Url, result.Level, result.VulName, result.Matched)
				} else {
					NotificationMsg += fmt.Sprintf("%v-[%v]-[%v]\n", result.Level, result.VulName, result.Matched)
				}
			}
		} else {
			if result.Url != "" {
				NotificationMsg += fmt.Sprintf("%v-[%v]-[%v]-[%v]\n", result.Url, result.Level, result.VulName, result.Matched)
			} else {
				NotificationMsg += fmt.Sprintf("%v-[%v]-[%v]\n", result.Level, result.VulName, result.Matched)
			}
		}
		notification.NotificationQueues["VulnerabilityScan"].Queue <- NotificationMsg
	}
	ResultQueues["VulnerabilityScan"].Queue <- interfaceSlice
}

func (h *handler) PageMonitoringInsert(result *types.PageMonit) {
	var interfaceSlice interface{}
	rootDomain, err := utils.Tools.GetRootDomain(result.Url)
	if err != nil {
		logger.SlogInfoLocal(fmt.Sprintf("%v GetRootDomain error: %v", result.Url, err))
	}
	result.RootDomain = rootDomain
	if result.Time == "" {
		result.Time = utils.Tools.GetTimeNow()
	}
	result.Project = h.GetAssetProject(rootDomain)
	interfaceSlice = &result
	ResultQueues["PageMonitoring"].Queue <- interfaceSlice
}

func (h *handler) PageMonitoringInsertBody(result *types.PageMonitBody) {
	var interfaceSlice interface{}
	interfaceSlice = &result
	ResultQueues["PageMonitoringBody"].Queue <- interfaceSlice
}

func (h *handler) RootDomain(result *types.RootDomain) {
	selector := bson.M{"domain": result.Domain}
	result.Project = h.GetAssetProject(result.Domain)
	var resultEx types.RootDomain
	err := mongodb.MongodbClient.FindOne("RootDomain", bson.M{"domain": result.Domain}, nil, &resultEx)
	tmpData := bson.M{
		"icp":      result.ICP,
		"tags":     result.Tags,
		"taskName": result.TaskName,
		"project":  result.Project,
		"time":     result.Time,
		"company":  result.Company,
	}
	if err != nil {
		// 出现错误表示 mongodb中不存在， 不存在则不进行处理 直接更新插入
	} else {
		// mongodb中存在
		// 查询mongodb中是否存在该domain 如果存在 则判断icp、company、project是否一样 如果一样就不需要更新，如果不一样需要更新，并且不一样的时候新的icp、company、project需要不为空
		if result.ICP == resultEx.ICP && result.Company == resultEx.Company && result.Project == resultEx.Project {
			return
		}
		if result.ICP == "" {
			tmpData["icp"] = resultEx.ICP
		}
		if result.Company == "" {
			tmpData["company"] = resultEx.Company
		}
		if result.Project == "" {
			tmpData["project"] = resultEx.Project
		}
	}
	update := bson.M{
		"$set": tmpData,
	}
	op := types.BulkUpdateOperation{
		Selector: selector,
		Update:   update,
	}
	ResultQueues["RootDomain"].Queue <- op
}

func (h *handler) APP(result *types.APP) {
	var interfaceSlice interface{}
	result.Project = h.GetAssetProject(result.ICP)
	interfaceSlice = &result
	ResultQueues["APP"].Queue <- interfaceSlice
}

func (h *handler) MP(result *types.MP) {
	var interfaceSlice interface{}
	result.Project = h.GetAssetProject(result.ICP)
	interfaceSlice = &result
	ResultQueues["MP"].Queue <- interfaceSlice
}
