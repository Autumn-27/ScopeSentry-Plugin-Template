// main-------------------------------------
// @file      : plugin_cmd.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/11/16 17:54
// -------------------------------------------

package main

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/plugins"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"path/filepath"
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

func main() {
	// plugin id
	plgId := utils.Tools.GenerateRandomString(8)
	// plugin module name
	plgModule := "AssetHandle"
	// plugin path
	plgPath := filepath.Join("D:\\code\\ScopeSentry\\ScopeSentry-Plugin-Template\\plugin\\AssetHandle\\plugin.go")

	plugin, err := plugins.LoadCustomPlugin(plgPath, plgModule, plgId)
	if err != nil {
		return
	}

	fmt.Printf("plugin name: %v", plugin.GetName())
	fmt.Printf("plugin module: %v", plugin.GetModule())
	fmt.Printf("plugin id: %v", plugin.GetPluginId())
	//op := options.PluginOption{
	//	Name:      plugin.GetName(),
	//	Module:    plugin.GetModule(),
	//	Parameter: plugin.GetParameter(),
	//	PluginId:  plugin.GetPluginId(),
	//	Result:    ,
	//	Custom:    p.Custom,
	//	TaskId:    p.TaskId,
	//}
}
