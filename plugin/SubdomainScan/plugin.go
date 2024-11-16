package plugin

import (
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
)

func GetName() string {
	return "SubdomainScan"
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

// Execute is a function that accepts a pointer to string as input.
// 执行的函数接受一个字符串作为输入。

// Parameters:
// - input: An empty interface (interface{}) that can hold any type, expected to be a string.
// - op: an options.PluginOption that contains configuration options for the operation.
// 参数：
// - input: 一个空接口（interface{}），可以容纳任何类型，预计是字符串。
// - op: 一个 options.PluginOption 类型，包含操作所需的配置选项。

// Execute Returns:
// - interface{}: the modified input or result of the operation.
// - error: an error if the operation fails, or nil if successful.
// 返回值：
// - interface{}: 无
// - error: 如果操作失败，返回错误；如果成功，返回 nil。
// result:
//
//		将结果发送到op.Result，结果类型为types.SubdomainResult
//	  Send the result to the op.result channel, the result type is types.SubdomainResult
func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	// 根据目标生成子域名，将子域名发送到result
	data, ok := input.(string)
	if !ok {
		// 说明不是http的资产，直接返回
		return nil, nil
	}

	result := make(chan string)
	go utils.Tools.ExecuteCommandToChan("whoami", []string{}, result)
	for i := range result {
		subdomainResult := types.SubdomainResult{
			Host:  i + "." + data,
			Type:  "A",
			Value: []string{"api.example.com", "api.example.net", "api.example.org"},
			IP:    []string{"192.168.100.10", "192.168.100.11", "192.168.100.12"},
			Time:  "",
			Tags:  []string{"demo-subdoamin"},
		}
		op.ResultFunc(subdomainResult)
	}
	return nil, nil
}
