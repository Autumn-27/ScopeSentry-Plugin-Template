// httpx-------------------------------------
// @file      : httpx.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/9/28 15:12
// -------------------------------------------

package httpx

import (
	"errors"
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/contextmanager"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/interfaces"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
	"strconv"
)

type Plugin struct {
	Name      string
	Module    string
	Parameter string
	PluginId  string
	Result    chan interface{}
	Custom    interface{}
	TaskId    string
	TaskName  string
}

func NewPlugin() *Plugin {
	return &Plugin{
		Name:     "httpx",
		Module:   "AssetMapping",
		PluginId: "3a0d994a12305cb15a5cb7104d819623",
	}
}

func (p *Plugin) SetTaskName(name string) {
	p.TaskName = name
}

func (p *Plugin) GetTaskName() string {
	return p.TaskName
}

func (p *Plugin) SetTaskId(id string) {
	p.TaskId = id
}

func (p *Plugin) GetTaskId() string {
	return p.TaskId
}
func (p *Plugin) SetCustom(cu interface{}) {
	p.Custom = cu
}

func (p *Plugin) GetCustom() interface{} {
	return p.Custom
}
func (p *Plugin) SetPluginId(id string) {
	p.PluginId = id
}

func (p *Plugin) GetPluginId() string {
	return p.PluginId
}

func (p *Plugin) SetResult(ch chan interface{}) {
	p.Result = ch
}

func (p *Plugin) SetName(name string) {
	p.Name = name
}

func (p *Plugin) GetName() string {
	return p.Name
}

func (p *Plugin) SetModule(module string) {
	p.Module = module
}

func (p *Plugin) GetModule() string {
	return p.Module
}

func (p *Plugin) Install() error {
	return nil
}

func (p *Plugin) Check() error {
	return nil
}
func (p *Plugin) UnInstall() error {
	return nil
}
func (p *Plugin) SetParameter(args string) {
	p.Parameter = args
}

func (p *Plugin) GetParameter() string {
	return p.Parameter
}

func (p *Plugin) Log(msg string, tp ...string) {
	var logTp string
	if len(tp) > 0 {
		logTp = tp[0] // 使用传入的参数
	} else {
		logTp = "i"
	}
	logger.PluginsLog(fmt.Sprintf("[Plugins %v] %v", p.GetName(), msg), logTp, p.GetModule(), p.GetPluginId())
}

func (p *Plugin) Execute(input interface{}) (interface{}, error) {
	data, ok := input.([]interface{})
	if !ok {
		//logger.SlogError(fmt.Sprintf("%v error: %v input is not types.AssetOther\n", p.Name, input))
		return nil, errors.New("input is not types.AssetOther")
	}
	var targetList []string
	for _, assetinterface := range data {
		asset, ok := assetinterface.(types.AssetOther)
		if !ok {
			//p.Log(fmt.Sprintf("assetinterface not types.AssetOther: %v", assetinterface), "w")
			continue
		}
		if asset.Type != "http" {
			p.Result <- asset
		} else {
			var url string
			if asset.Port != "" {
				url = asset.Host + ":" + asset.Port + asset.UrlPath
			} else {
				url = asset.Host + asset.UrlPath
			}
			targetList = append(targetList, url)
		}
	}
	parameter := p.GetParameter()
	cdncheck := "false"
	screenshot := false
	tlsprobe := true
	FollowRedirects := true
	bypassHeader := false
	screenshotTimeout := 10
	executionTimeout := 10
	if parameter != "" {
		args, err := utils.Tools.ParseArgs(parameter, "cdncheck", "screenshot", "st", "tlsprobe", "fr", "et", "bh")
		if err != nil {
		} else {
			for key, value := range args {
				if value != "" {
					switch key {
					case "cdncheck":
						cdncheck = value
					case "screenshot":
						if value == "true" {
							screenshot = true
						}
					case "tlsprobe":
						if value == "false" {
							tlsprobe = false
						}
					case "st":
						screenshotTimeout, _ = strconv.Atoi(value)
					case "fr":
						if value == "false" {
							FollowRedirects = false
						}
					case "et":
						executionTimeout, _ = strconv.Atoi(value)
					case "bh":
						if value == "true" {
							bypassHeader = true
						}
					default:
						continue
					}
				}
			}
		}
	}
	httpxResultsHandler := func(r types.AssetHttp) {
		p.Result <- r
	}

<<<<<<< HEAD
	executionTimeout := 10
	bypassHeader := false
=======
>>>>>>> 166404a45d573ee1202c4c05526269963c677b69
	utils.Requests.Httpx(targetList, httpxResultsHandler, cdncheck, screenshot, screenshotTimeout, tlsprobe, FollowRedirects, contextmanager.GlobalContextManagers.GetContext(p.GetTaskId()), executionTimeout, bypassHeader)
	return nil, nil
}

func (p *Plugin) Clone() interfaces.Plugin {
	return &Plugin{
		Name:     p.Name,
		Module:   p.Module,
		PluginId: p.PluginId,
		Custom:   p.Custom,
		TaskId:   p.TaskId,
	}
}
