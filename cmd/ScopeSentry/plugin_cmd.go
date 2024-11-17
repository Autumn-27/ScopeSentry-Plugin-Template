// main-------------------------------------
// @file      : plugin_cmd.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/11/16 17:54
// -------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/bigcache"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/config"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/configupdater"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/contextmanager"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/global"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/handler"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/mongodb"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/notification"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/pebbledb"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/plugins"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/pool"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/redis"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/results"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

// All modules

//"TargetHandler"
//"SubdomainScan"
//"SubdomainSecurity"
//"AssetMapping"
//"PortScanPreparation"
//"PortScan"
//"PortFingerprint"
//"AssetHandle"
//"URLScan"
//"URLSecurity"
//"WebCrawler"
//"DirScan"
//"VulnerabilityScan"

// demo data
var (
	AssetHttpData = types.AssetHttp{
		Time:          "2024-11-16T08:30:00Z",
		LastScanTime:  "2024-11-15T10:00:00Z",
		TLSData:       nil, // 如果没有 TLS 数据可以设置为 nil
		Hashes:        map[string]interface{}{"sha256": "abcdef1234567890"},
		CDNName:       "Cloudflare",
		Port:          "443",
		URL:           "https://example.com",
		Title:         "Example Site",
		Type:          "web",
		Error:         "",
		ResponseBody:  "Example response body here.",
		Host:          "example.com",
		IP:            "192.168.1.1",
		Screenshot:    "path/to/screenshot.png",
		FavIconMMH3:   "abc123",
		FaviconPath:   "/assets/favicon.ico",
		RawHeaders:    "HTTP/1.1 200 OK\nContent-Type: text/html\n",
		Jarm:          "JARM hash data here",
		Technologies:  []string{"NGINX", "React", "Node.js"},
		StatusCode:    200,
		ContentLength: 1234,
		CDN:           true,
		Webcheck:      true,
		Project:       "Project X",
		IconContent:   "base64iconcontenthere",
		Domain:        "example.com",
		TaskName:      []string{"Task 1", "Task 2"},
		WebServer:     "nginx/1.21.0",
		Service:       "Web Hosting",
		RootDomain:    "example.com",
		Tags:          []string{"production", "ssl", "security"},
	}
	AssetOtherData = types.AssetOther{
		Time:         "2024-11-16T08:30:00Z",
		LastScanTime: "2024-11-15T10:00:00Z",
		Host:         "other-example.com",
		IP:           "192.168.1.2",
		Port:         "8080",
		Service:      "HTTP Server",
		TLS:          true,
		Transport:    "TCP",
		Version:      "1.0.0",
		Raw:          json.RawMessage(`{"metadataField": "value"}`),
		Project:      "Project Y",
		Type:         "service",
		Tags:         []string{"internal", "ssl", "test"},
		TaskName:     []string{"Task A", "Task B"},
		RootDomain:   "other-example.com",
	}
)

func main() {
	global.DatabaseEnabled = false
	Init()
	global.AppConfig.Debug = false
	_, filePath, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("无法获取当前文件路径")
	}
	parentDir := filepath.Dir(filePath)
	plgPath := filepath.Join(parentDir, "..", "..", "plugin")

	TestAssetHandle(plgPath)
	TestSubdomainScan(plgPath)
}

func TestAssetHandle(plgPath string) {
	// plugin id
	plgId := utils.Tools.GenerateRandomString(8)
	// plugin module name
	plgModule := "AssetHandle"
	// plugin path
	plgPath = filepath.Join(plgPath, "AssetHandle", "plugin.go")

	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
	if err != nil {
		return
	}

	fmt.Printf("plugin name: %v\n", plugin.GetName())
	fmt.Printf("plugin module: %v\n", plugin.GetModule())
	fmt.Printf("plugin id: %v\n", plugin.GetPluginId())
	result := make(chan interface{})
	plugin.SetParameter("")
	plugin.SetTaskId("1111")
	plugin.SetTaskName("demo")
	plugin.SetResult(result)
	fmt.Printf("AssetHttpData original tags: %v\n", AssetHttpData.Tags)
	_, err = plugin.Execute(&AssetHttpData)
	if err != nil {
		return
	}
	fmt.Printf("AssetHttpData tags%v\n", AssetHttpData.Tags)
}

func TestSubdomainScan(plgPath string) {
	// plugin id
	plgId := utils.Tools.GenerateRandomString(8)
	// plugin module name
	plgModule := "SubdomainScan"
	// plugin path
	plgPath = filepath.Join(plgPath, "SubdomainScan", "plugin.go")

	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("plugin name: %v\n", plugin.GetName())
	fmt.Printf("plugin module: %v\n", plugin.GetModule())
	fmt.Printf("plugin id: %v\n", plugin.GetPluginId())
	result := make(chan interface{})
	plugin.SetParameter("")
	plugin.SetTaskId("1111")
	plugin.SetTaskName("demo")
	plugin.SetResult(result)
	go func() {
		for r := range result {
			jsonData, _ := json.Marshal(r)
			fmt.Printf("result %v", string(jsonData))
		}
	}()
	_, err = plugin.Execute("example.com")
	if err != nil {
		return
	}
	time.Sleep(3 * time.Second)
	fmt.Printf("plugin name: %v\n", plugin.GetName())
	fmt.Printf("plugin module: %v\n", plugin.GetModule())
	fmt.Printf("plugin id: %v\n", plugin.GetPluginId())

}

func Init() {
	// 初始化系统信息
	config.Initialize()
	global.VERSION = "1.5"
	var err error
	if global.DatabaseEnabled {
		// 初始化mongodb连接
		mongodb.Initialize()
		// 初始化redis连接
		redis.Initialize()
	}
	// 初始化日志模块
	err = logger.NewLogger()
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	// 初始化任务计数器
	handler.InitHandle()
	// 更新配置、加载字典
	if global.DatabaseEnabled {
		configupdater.Initialize()
	}
	// 初始化模块配置
	err = config.ModulesInitialize()
	if err != nil {
		log.Fatalf("Failed to init ModulesConfig: %v", err)
		return
	}
	// 初始化上下文管理器
	contextmanager.NewContextManager()
	// 初始化tools
	utils.InitializeTools()
	utils.InitializeDnsTools()
	utils.InitializeRequests()
	utils.InitializeResults()
	// 初始化通知模块
	notification.InitializeNotification()
	// 初始化协程池
	pool.Initialize()
	// 初始化个模块的协程池
	pool.PoolManage.InitializeModulesPools(config.ModulesConfig)
	go pool.StartMonitoring()
	// 初始化内存缓存
	err = bigcache.Initialize()
	if err != nil {
		logger.SlogErrorLocal(fmt.Sprintf("bigcache Initialize error: %v", err))
		return
	}
	// 初始化本地持久化缓存
	pebbledbSetting := pebbledb.Settings{
		DBPath:       filepath.Join(global.AbsolutePath, "data", "pebbledb"),
		CacheSize:    64 << 20,
		MaxOpenFiles: 500,
	}
	pebbledbOption := pebbledb.GetPebbleOptions(&pebbledbSetting)
	if !global.AppConfig.Debug {
		pebbledbOption.Logger = nil
	}
	pedb, err := pebbledb.NewPebbleDB(pebbledbOption, pebbledbSetting.DBPath)
	if err != nil {
		return
	}
	pebbledb.PebbleStore = pedb
	defer func(PebbleStore *pebbledb.PebbleDB) {
		_ = PebbleStore.Close()
	}(pebbledb.PebbleStore)

	// 初始化结果处理队列，(正常的时候将该初始化放入任务开始时，任务执行完毕关闭结果队列)
	results.InitializeResultQueue()
	defer results.Close()

	// 初始化全局插件管理器
	//plugins.GlobalPluginManager = plugins.NewPluginManager()
	//err = plugins.GlobalPluginManager.InitializePlugins()
	//if err != nil {
	//	log.Fatalf("Failed to init plugins: %v", err)
	//	return
	//}

}