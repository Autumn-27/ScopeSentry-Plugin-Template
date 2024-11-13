// main-------------------------------------
// @file      : restructure_cmd.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/9/6 22:00
// -------------------------------------------

package main

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/bigcache"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/config"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/configupdater"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/contextmanager"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/global"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/handler"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/mongodb"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/node"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/notification"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/pebbledb"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/plugins"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/pool"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/redis"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/results"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/task"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func main() {
	// 初始化系统信息
	config.Initialize()
	global.VERSION = "1.5"
	global.AppConfig.Debug = true
	var err error
	// 初始化mongodb连接
	mongodb.Initialize()
	// 初始化redis连接
	redis.Initialize()
	// 初始化日志模块
	err = logger.NewLogger()
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	// 初始化任务计数器
	handler.InitHandle()
	// 更新配置、加载字典
	configupdater.Initialize()
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
	plugins.GlobalPluginManager = plugins.NewPluginManager()
	err = plugins.GlobalPluginManager.InitializePlugins()
	if err != nil {
		log.Fatalf("Failed to init plugins: %v", err)
		return
	}
	// 性能监控
	go pprof()
	//taskE := options.TaskOptions{
	//	ID:       "1",
	//	TaskName: "test",
	//	//SubdomainScan:       []string{"subfinder", "ksubdomain"},
	//	//SubdomainSecurity:   []string{"SubdomainTakeover"},
	//	//PortScanPreparation: []string{"SkipCdn"},
	//	PortScan: []string{"RustScan"},
	//	//PortFingerprint:     []string{"fingerprintx"},
	//	AssetMapping: []string{"httpx"},
	//	//AssetHandle:         []string{"WebFingerprint"},
	//	URLScan:     []string{"katana", "wayback"},
	//	URLSecurity: []string{"sensitive"},
	//	//WebCrawler:          []string{"rad"},
	//	//DirScan:             []string{"SentryDir"},
	//	//VulnerabilityScan: []string{"nuclei"},
	//	SubdomainFilename: "66dda6ee3687eb004e6b3bda",
	//	Parameters: map[string]map[string]string{
	//		"SubdomainScan": {
	//			"ksubdomain": "-subfile 66dda6ee3687eb004e6b3bda",
	//		},
	//		"DirScan": {
	//			"SentryDir": "-d 6712604bcad762b691bf0d93",
	//		},
	//		"VulnerabilityScan": {
	//			"nuclei": "-t *",
	//		},
	//	},
	//}
	//jsonStr, err := utils.Tools.StructToJSON(taskE)
	//if err != nil {
	//	return
	//}
	//pebbledb.PebbleStore.Put([]byte("task:1"), []byte(jsonStr))
	//
	//pebbledb.PebbleStore.Put([]byte("1:http://39.105.160.88:83"), []byte("1"))

	//taskE = options.TaskOptions{
	//	ID:       "2",
	//	TaskName: "test",
	//	//SubdomainScan:       []string{"subfinder", "ksubdomain"},
	//	//SubdomainSecurity:   []string{"SubdomainTakeover"},
	//	//PortScanPreparation: []string{"SkipCdn"},
	//	PortScan: []string{"RustScan"},
	//	//PortFingerprint:     []string{"fingerprintx"},
	//	//AssetMapping:        []string{"httpx"},
	//	//AssetHandle:         []string{"WebFingerprint"},
	//	//URLScan:             []string{"katana", "wayback"},
	//	//URLSecurity:         []string{"sensitive"},
	//	//WebCrawler:          []string{"rad"},
	//	//DirScan:             []string{"SentryDir"},
	//	VulnerabilityScan: []string{"nuclei"},
	//	SubdomainFilename: "66dda6ee3687eb004e6b3bda",
	//	Parameters: map[string]map[string]string{
	//		"SubdomainScan": {
	//			"ksubdomain": "-subfile 66dda6ee3687eb004e6b3bda",
	//		},
	//		"DirScan": {
	//			"SentryDir": "-d 6712604bcad762b691bf0d93",
	//		},
	//		"VulnerabilityScan": {
	//			"nuclei": "-t 66dd6b0edcf0ee4a4fef8bf7.yaml",
	//		},
	//	},
	//}
	//jsonStr, err = utils.Tools.StructToJSON(taskE)
	//if err != nil {
	//	return
	//}
	//pebbledb.PebbleStore.Put([]byte("task:2"), []byte(jsonStr))
	//pebbledb.PebbleStore.Put([]byte("2:http://123.56.16.123"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("1:http://39.105.160.88:666"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("1:http://39.105.160.88:80"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("1:http://39.105.160.88:8081"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("1:https://baidu.com"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("1:https://hackerone.com"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("2:baidu.com"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("2:google.com"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("2:tes1t.com"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("2:tes2t.com"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("2:tes3t.com"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("2:tes4t.com"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("2:tes5t.com"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("2:tes6t.com"), []byte("1"))
	//pebbledb.PebbleStore.Put([]byte("2:tes7t.com"), []byte("1"))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done() // 减少计数器，表示任务完成
		task.GetTask()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		node.Register()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		go configupdater.RefreshConfig()
	}()
	time.Sleep(10 * time.Second)
	wg.Wait()
}

func pprof() {
	if global.AppConfig.Debug {
		go func() {
			_ = http.ListenAndServe("0.0.0.0:6060", nil)
		}()
		//go DebugMem()
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			fmt.Println("收到终止信号:", sig)
		}()
	}
}
