package plugin

import (
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
)

func GetName() string {
	return "AssetHandleDemo"
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

// Execute is a function that accepts a pointer to either types.AssetHttp or AssetOther as input.
// The input can be modified directly, which will change the result accordingly.
// 执行的函数接受一个指向 types.AssetHttp 或 AssetOther 的指针作为输入。
// 输入可以直接修改，从而改变结果。

// Parameters:
// - input: an interface{} that can hold any type, but it is expected to be a pointer to either types.AssetHttp or AssetOther.
// - op: an options.PluginOption that contains configuration options for the operation.
// 参数：
// - input: 一个空接口（interface{}），可以容纳任何类型，预计是指向 types.AssetHttp 或 AssetOther 的指针。
// - op: 一个 options.PluginOption 类型，包含操作所需的配置选项。

// Execute Returns:
// - interface{}: the modified input or result of the operation.
// - error: an error if the operation fails, or nil if successful.
// 返回值：
// - interface{}: 无
// - error: 如果操作失败，返回错误；如果成功，返回 nil。
func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	// 给http服务打tag
	// Tag the http service
	httpResult, ok := input.(*types.AssetHttp)
	if !ok {
		// 说明不是http的资产，直接返回
		return nil, nil
	}
	httpResult.Tags = append(httpResult.Tags, "Demo")
	return nil, nil
}
