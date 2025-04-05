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
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/pebbledb"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/plugins"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/pool"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/redis"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/results"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	plugin "github.com/Autumn-27/ScopeSentry-Scan/plugin/URLSecurity/DependencyConfusion"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
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
		Title:         "",
		Type:          "web",
		Error:         "",
		ResponseBody:  "Example response body style1/css/ListRange.css 主账套 login.jsp",
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
	UrlRes = types.UrlResult{
		Output:   "http://39.105.160.88:666/djwaklg.php",
		ResultId: "resultId",
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
	fmt.Println(plgPath)

	//TestAssetHandle(plgPath)
	//TestSubdomainScan(plgPath)
	//TestEHole(plgPath)
	//TestEHoleDebug()
	//TestSSRFScanDebug()
	//TestSSRFScan(plgPath)
	//TestFofa(plgPath)
	//TestFofaDebug()
	//TestXray(plgPath)
	//TestAi(plgPath)
	//TestDep(plgPath)
	TestDepDebug()
}

func TestDepDebug() {
	op := options.PluginOption{
		Name:      "DependencyConfusion",
		Module:    "URLSecurity",
		Parameter: "",
		PluginId:  "",
		Ctx:       contextmanager.GlobalContextManagers.GetContext("111111"),
	}
	input := types.UrlResult{
		Output: "https://cxxxxxx/js/cxxx5.js",
		Ext:    ".js",
		Status: 200,
		Body: `(this._delta+=i,this._lastMousePos=this._map.mouseEventToContainerPoint(t),this._startTime||(this._startTime=+new Date),Math.max(e-(+new Date-this._startTime),0));clearTimeout(this._timer),this._timer=setTimeout(a(this._performZoom,this),i),Ri(t)},_performZoom:function(){var t=this._map,i=t.getZoom(),e=this._map.options.zoomSnap||0,n=(t._stop(),this._delta/(4*this._map.options.wheelPxPerZoomLevel)),n=4*Math.log(2/(1+Math.exp(-Math.abs(n))))/Math.LN2,e=e?Math.ceil(n/e)*e:n,n=t._limitZoom(i+(0<this._delta?e:-e))-i;this._delta=0,this._startTime=null,n&&("center"===t.options.scrollWheelZoom?t.setZoom(i+n):t.setZoomAround(this._lastMousePos,i+n))}})),Et=(A.addInitHook("addHandler","scrollWheelZoom",kt),A.mergeOptions({tapHold:P.touchNative&&P.safari&&P.mobile,tapTolerance:15}),n.extend({addHooks:function(){S(this._map._container,"touchstart",this._onDown,this)},removeHooks:function(){E(this._map._container,"touchstart",this._onDown,this)},_onDown:function(t){var i;clearTimeout(this._holdTimeout),1===t.touches.length&&(i=t.touches[0],this._startPos=this._newPos=new p(i.clientX,i.clientY),this._holdTimeout=setTimeout(a(function(){this._cancel(),this._isTapValid()&&(S(document,"touchend",B),S(document,"touchend touchcancel",this._cancelClickPrevent),this._simulateEvent("contextmenu",i))},this),600),S(document,"touchend touchcancel contextmenu",this._cancel,this),S(document,"touchmove",this._onMove,this))},_cancelClickPrevent:function t(){E(document,"touchend",B),E(document,"touchend touchcancel",t)},_cancel:function(){clearTimeout(this._holdTimeout),E(document,"touchend touchcancel contextmenu",this._cancel,this),E(document,"touchmove",this._onMove,this)},_onMove:function(t){t=t.touches[0];this._newPos=new p(t.clientX,t.clientY)},_isTapValid:function(){return this._newPos.distanceTo(this._startPos)<=this._map.options.tapTolerance},_simulateEvent:function(t,i){t=new MouseEvent(t,{bubbles:!0,cancelable:!0,view:window,screenX:i.screenX,screenY:i.screenY,clientX:i.clientX,clientY:i.clientY});t._simulated=!0,i.target.dispatchEvent(t)}})),Bt=(A.addInitHook("addHandler","tapHold",Et),A.mergeOptions({touchZoom:P.touch,bounceAtZoomLimits:!0}),n.extend({addHooks:function(){z(this._map._container,"leaflet-touch-zoom"),S(this._map._container,"touchstart",this._onTouchStart,this)},removeHooks:function(){M(this._map._container,"leaflet-touch-zoom"),E(this._map._container,"touchstart",this._onTouchStart,this)},_onTouchStart:function(t){var i,e,n=this._map;!t.touches||2!==t.touches.length||n._animatingZoom||this._zooming||(i=n.mouseEventToContainerPoint(t.touches[0]),e=n.mouseEventToContainerPoint(t.touches[1]),this._centerPoint=n.getSize()._divideBy(2),this._startLatLng=n.containerPointToLatLng(this._centerPoint),"center"!==n.options.touchZoom&&(this._pinchStartLatLng=n.containerPointToLatLng(i.add(e)._divideBy(2))),this._startDist=i.distanceTo(e),this._startZoom=n.getZoom(),this._moved=!1,this._zooming=!0,n._stop(),S(document,"touchmove",this._onTouchMove,this),S(document,"touchend touchcancel",this._onTouchEnd,this),B(t))},_onTouchMove:function(t){if(t.touches&&2===t.touches.length&&this._zooming){var i=this._map,e=i.mouseEventToContainerPoint(t.touches[0]),n=i.mouseEventToContainerPoint(t.touches[1]),o=e.distanceTo(n)/this._startDist;if(this._zoom=i.getScaleZoom(o,this._startZoom),!i.options.bounceAtZoomLimits&&(this._zoom<i.getMinZoom()&&o<1||this._zoom>i.getMaxZoom()&&1<o)&&(this._zoom=i._limitZoom(this._zoom)),"center"===i.options.touchZoom){if(this._center=this._startLatLng,1==o)return}else{e=e._add(n)._divideBy(2)._subtract(this._centerPoint);if(1==o&&0===e.x&&0===e.y)return;this._center=i.unproject(i.project(this._pinchStartLatLng,this._zoom).subtract(e),this._zoom)}this._moved||(i._moveStart(!0,!1),this._moved=!0),r(this._animRequest);n=a(i._move,i,this._center,this._zoom,{pinch:!0,round:!1});this._animRequest=x(n,this,!0),B(t)}},_onTouchEnd:function(){this._moved&&this._zooming?(this._zooming=!1,r(this._animRequest),E(document,"touchmove",this._onTouchMove,this),E(document,"touchend touchcancel",this._onTouchEnd,this),this._map.options.zoomAnimation?this._map._animateZoom(this._center,this._map._limitZoom(this._zoom),!0,this._map.options.zoomSnap):this._map._resetView(this._center,this._map._limitZoom(this._zoom))):this._zooming=!1}})),qe=(A.addInitHook("addHandler","touchZoom",Bt),A.BoxZoom=_t,A.DoubleClickZoom=Ct,A.Drag=Zt,A.Keyboard=St,A.ScrollWheelZoom=kt,A.TapHold=Et,A.TouchZoom=Bt,t.Bounds=m,t.Browser=P,t.CRS=ot,t.Canvas=De,t.Circle=me,t.CircleMarker=pe,t.Class=it,t.Control=I,t.DivIcon=Be,t.DivOverlay=O,t.DomEvent=mt,t.DomUtil=pt,t.Draggable=Xi,t.Evented=et,t.FeatureGroup=he,t.GeoJSON=ve,t.GridLayer=Ae,t.Handler=n,t.Icon=le,t.ImageOverlay=Ce,t.LatLng=v,t.LatLngBounds=s,t.Layer=o,t.LayerGroup=ae,t.LineUtil=gt,t.Map=A,t.Marker=de,t.Mixin=ft,t.Path=_e,t.Point=p,t.PolyUtil=vt,t.Polygon=ge,t.Polyline=fe,t.Popup=ke,t.PosAnimation=Wi,t.Projection=wt,t.Rectangle=Ve,t.Renderer=Ne,t.SVG=Fe,t.SVGOverlay=Se,t.TileLayer=Ie,t.Tooltip=Ee,t.Transformation=at,t.Util=tt,t.VideoOverlay=Ze,t.bind=a,t.bounds=f,t.canvas=je,t.circle=function(t,i,e){return new me(t,i,e)},t.circleMarker=function(t,i){return new pe(t,i)},t.control=Fi,t.divIcon=function(t){return new Be(t)},t.extend=l,t.featureGroup=function(t,i){return new he(t,i)},t.geoJSON=Me,t.geoJson=zt,t.gridLayer=function(t){return new Ae(t)},t.icon=function(t){return new le(t)},t.imageOverlay=function(t,i,e){return new Ce(t,i,e)},t.latLng=w,t.latLngBounds=g,t.layerGroup=function(t,i){return new ae(t,i)},t.map=function(t,i){return new A(t,i)},t.marker=function(t,i){return new de(t,i)},t.point=_,t.polygon=function(t,i){return new ge(t,i)},t.polyline=function(t,i){return new fe(t,i)},t.popup=function(t,i){return new ke(t,i)},t.rectangle=function(t,i){return new Ve(t,i)},t.setOptions=c,t.stamp=h,t.svg=Ue,t.svgOverlay=function(t,i,e){return new Se(t,i,e)},t.tileLayer=Oe,t.tooltip=function(t,i){return new Ee(t,i)},t.transformation=ht,t.version="1.8.0",t.videoOverlay=function(t,i,e){return new Ze(t,i,e)},window.L);t.noConflict=function(){return window.L=qe,this},window.L=t});
//# sourceMappingURL=leaflet.js.map`,
	}
	data, err := ioutil.ReadFile("C:\\Users\\autumn\\Desktop\\a.txt")
	if err != nil {
		log.Fatal(err)
	}
	input.Body = string(data)
	plugin.Execute(input, op)
}

func TestDep(plgPath string) {
	// plugin id
	plgId := utils.Tools.GenerateRandomString(8)
	// plugin module name
	plgModule := "URLSecurity"
	// plugin path
	plgPath = filepath.Join(plgPath, "URLSecurity", "DependencyConfusion", "plugin.go")
	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
	if err != nil {
		fmt.Println(err)
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
	plugin.Install()
	var input []types.AssetHttp
	input = append(input, AssetHttpData)
	go func() {
		_, err = plugin.Execute(input)
		if err != nil {
			return
		}
	}()
	for data := range result {
		fmt.Println(data) // 打印接收到的数据
	}
}

func TestAi(plgPath string) {
	// plugin id
	plgId := utils.Tools.GenerateRandomString(8)
	// plugin module name
	plgModule := "VulnerabilityScan"
	// plugin path
	plgPath = filepath.Join(plgPath, "VulnerabilityScan", "AI-Infra-Guard", "plugin.go")

	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
	if err != nil {
		fmt.Println(err)
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
	plugin.Install()
	var input []types.AssetHttp
	input = append(input, AssetHttpData)
	go func() {
		_, err = plugin.Execute(input)
		if err != nil {
			return
		}
	}()
	for data := range result {
		fmt.Println(data) // 打印接收到的数据
	}
}

func TestXray(plgPath string) {
	// plugin id
	plgId := utils.Tools.GenerateRandomString(8)
	// plugin module name
	plgModule := "PassiveScan"
	// plugin path
	plgPath = filepath.Join(plgPath, "PassiveScan", "xray", "plugin.go")

	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
	if err != nil {
		fmt.Println(err)
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
	plugin.Install()
	go func() {
		_, err = plugin.Execute("baidu.com")
		if err != nil {
			return
		}
	}()
	var input string

	fmt.Print("请输入一个值: ")
	fmt.Scanln(&input) // 获取用户输入
	if input == "2" {
		plugin.SetCustom("close task")
	}
	for data := range result {
		fmt.Println(data) // 打印接收到的数据
	}
}

//
//func TestFofaDebug() {
//	plg := customplugin.NewPlugin("TargetHandler", "11111", plugin.Install, plugin.Check, plugin.Execute, plugin.Uninstall, plugin.GetName)
//	plg.SetParameter("")
//	plg.Execute("baidu.com")
//}
//
//func TestFofa(plgPath string) {
//	// plugin id
//	plgId := utils.Tools.GenerateRandomString(8)
//	// plugin module name
//	plgModule := "TargetHandler"
//	// plugin path
//	plgPath = filepath.Join(plgPath, "TargetHandler", "fofa", "plugin.go")
//
//	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
//	if err != nil {
//		return
//	}
//
//	fmt.Printf("plugin name: %v\n", plugin.GetName())
//	fmt.Printf("plugin module: %v\n", plugin.GetModule())
//	fmt.Printf("plugin id: %v\n", plugin.GetPluginId())
//	result := make(chan interface{})
//	plugin.SetParameter("")
//	plugin.SetTaskId("1111")
//	plugin.SetTaskName("demo")
//	plugin.SetResult(result)
//
//	go func() {
//		_, err = plugin.Execute("baidu.com")
//		if err != nil {
//			return
//		}
//	}()
//	for data := range result {
//		fmt.Println(data) // 打印接收到的数据
//	}
//}

//
//func TestSSRFScanDebug() {
//	plg := customplugin.NewPlugin("AssetHandle", "11111", plugin.Install, plugin.Check, plugin.Execute, plugin.Uninstall, plugin.GetName)
//	plg.SetParameter("-parfile 674c411aaa621e265dc1815a -dnslog fde390d9.log.dnslog.sbs")
//	plg.Execute(UrlRes)
//}
//
//func TestSSRFScan(plgPath string) {
//	// plugin id
//	plgId := utils.Tools.GenerateRandomString(8)
//	// plugin module name
//	plgModule := "URLSecurity"
//	// plugin path
//	plgPath = filepath.Join(plgPath, "URLSecurity", "SSRFScan", "plugin.go")
//
//	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
//	if err != nil {
//		return
//	}
//
//	fmt.Printf("plugin name: %v\n", plugin.GetName())
//	fmt.Printf("plugin module: %v\n", plugin.GetModule())
//	fmt.Printf("plugin id: %v\n", plugin.GetPluginId())
//	result := make(chan interface{})
//	plugin.SetParameter("-parfile 674c411aaa621e265dc1815a -dnslog fde390d9.log.dnslog.sbs")
//	plugin.SetTaskId("1111")
//	plugin.SetTaskName("demo")
//	plugin.SetResult(result)
//	_, err = plugin.Execute(UrlRes)
//	if err != nil {
//		return
//	}
//}
//
//func TestEHoleDebug() {
//	op := options.PluginOption{
//		Name:      "EHole",
//		Module:    "AssetHandle",
//		Parameter: "-finger dwa -thread 20",
//		PluginId:  "11111",
//		Ctx:       contextmanager.GlobalContextManagers.GetContext("111111"),
//	}
//	plugin.Execute(&AssetHttpData, op)
//}
//
//func TestEHole(plgPath string) {
//	// plugin id
//	plgId := utils.Tools.GenerateRandomString(8)
//	// plugin module name
//	plgModule := "AssetHandle"
//	// plugin path
//	plgPath = filepath.Join(plgPath, "AssetHandle", "ehole", "plugin.go")
//
//	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
//	if err != nil {
//		return
//	}
//
//	fmt.Printf("plugin name: %v\n", plugin.GetName())
//	fmt.Printf("plugin module: %v\n", plugin.GetModule())
//	fmt.Printf("plugin id: %v\n", plugin.GetPluginId())
//	result := make(chan interface{})
//	plugin.SetParameter("-finger dwa -thread 20")
//	plugin.SetTaskId("1111")
//	plugin.SetTaskName("demo")
//	plugin.SetResult(result)
//	fmt.Printf("AssetHttpData original Technologies: %v\n", AssetHttpData.Technologies)
//	_, err = plugin.Execute(&AssetHttpData)
//	if err != nil {
//		return
//	}
//	fmt.Printf("AssetHttpData Technologies%v\n", AssetHttpData.Technologies)
//}
//
//func TestAssetHandle(plgPath string) {
//	// plugin id
//	plgId := utils.Tools.GenerateRandomString(8)
//	// plugin module name
//	plgModule := "AssetHandle"
//	// plugin path
//	plgPath = filepath.Join(plgPath, "AssetHandle", "demo", "plugin.go")
//
//	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
//	if err != nil {
//		return
//	}
//
//	fmt.Printf("plugin name: %v\n", plugin.GetName())
//	fmt.Printf("plugin module: %v\n", plugin.GetModule())
//	fmt.Printf("plugin id: %v\n", plugin.GetPluginId())
//	result := make(chan interface{})
//	plugin.SetParameter("")
//	plugin.SetTaskId("1111")
//	plugin.SetTaskName("demo")
//	plugin.SetResult(result)
//	fmt.Printf("AssetHttpData original tags: %v\n", AssetHttpData.Tags)
//	_, err = plugin.Execute(&AssetHttpData)
//	if err != nil {
//		return
//	}
//	fmt.Printf("AssetHttpData tags%v\n", AssetHttpData.Tags)
//}
//
//func TestSubdomainScan(plgPath string) {
//	// plugin id
//	plgId := utils.Tools.GenerateRandomString(8)
//	// plugin module name
//	plgModule := "SubdomainScan"
//	// plugin path
//	plgPath = filepath.Join(plgPath, "SubdomainScan", "demo", "plugin.go")
//
//	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
//	if err != nil {
//		fmt.Printf("%v", err)
//		return
//	}
//
//	fmt.Printf("plugin name: %v\n", plugin.GetName())
//	fmt.Printf("plugin module: %v\n", plugin.GetModule())
//	fmt.Printf("plugin id: %v\n", plugin.GetPluginId())
//	result := make(chan interface{})
//	plugin.SetParameter("")
//	plugin.SetTaskId("1111")
//	plugin.SetTaskName("demo")
//	plugin.SetResult(result)
//	go func() {
//		for r := range result {
//			jsonData, _ := json.Marshal(r)
//			fmt.Printf("result %v", string(jsonData))
//		}
//	}()
//	_, err = plugin.Execute("example.com")
//	if err != nil {
//		return
//	}
//	time.Sleep(3 * time.Second)
//	fmt.Printf("plugin name: %v\n", plugin.GetName())
//	fmt.Printf("plugin module: %v\n", plugin.GetModule())
//	fmt.Printf("plugin id: %v\n", plugin.GetPluginId())
//
//}

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
	moduleConfig := config.ModulesConfigStruct{
		MaxGoroutineCount: 10,
	}
	utils.Tools.EnsureFilePathExists(config.ModulesConfigPath)
	// 写入模块配置
	err = utils.Tools.WriteYAMLFile(config.ModulesConfigPath, moduleConfig)
	//// 初始化模块配置
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
