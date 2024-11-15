// main-------------------------------------
// @file      : test_yaegi.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/10/17 19:04
// -------------------------------------------

package main

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/symbols"
	"github.com/Autumn-27/ScopeSentry-Scan/modules/customplugin"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"os"
	"path/filepath"
)

const src = `package foo

import "github.com/Autumn-27/ScopeSentry-Scan/pkg/logger"

func Bar(s string) string {
	logger.SlogInfoLocal("system config load begin")
	return s + "-Foo"
}
`

func main() {
	logger.NewLogger()
	// 获取可执行文件的目录
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execDir := filepath.Dir(execPath)
	fmt.Printf("Executable directory: %s\n", execDir)

	// 初始化 yaegi 解释器
	interp := interp.New(interp.Options{})

	// 加载标准库和符号
	interp.Use(stdlib.Symbols)
	interp.Use(symbols.Symbols)

	// 加载插件
	pluginPath := filepath.Join("D:\\code\\ScopeSentry\\ScopeSentry-Plugin-Template\\plugin", "plugin")
	fmt.Printf("Loading plugin from: %s\n", pluginPath) // 打印插件路径以确认
	_, err = interp.EvalPath(pluginPath)
	if err != nil {
		panic(err)
	}

	// 获取 foo.Bar 函数
	v, err := interp.Eval("plugin.Execute")
	if err != nil {
		panic(err)
	}
	// 将值转换为函数
	executeFunc := v.Interface().(func(input interface{}, op options.PluginOption) (interface{}, error))
	v, err = interp.Eval("plugin.GetName")
	if err != nil {
		panic(err)
	}
	getNameFunc := v.Interface().(func() string)

	v, err = interp.Eval("plugin.Install")
	if err != nil {
		panic(err)
	}
	installFunc := v.Interface().(func() error)

	v, err = interp.Eval("plugin.Check")
	if err != nil {
		panic(err)
	}
	checkFunc := v.Interface().(func() error)

	v, err = interp.Eval("plugin.Uninstall")
	if err != nil {
		panic(err)
	}
	uninstallFunc := v.Interface().(func() error)
	plg := customplugin.NewPlugin("test", "", installFunc, checkFunc, executeFunc, uninstallFunc, getNameFunc)
	nePlg := plg.Clone()
	fmt.Println(nePlg.Execute("d"))
}
