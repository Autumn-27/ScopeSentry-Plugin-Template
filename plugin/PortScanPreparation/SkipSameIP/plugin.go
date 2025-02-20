// SkipSameIP-------------------------------------
// @file      : plugin.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/2/20 20:47
// -------------------------------------------

package plugin

import (
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/results"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
)

func GetName() string {
	return "SkipSameIP"
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

func Execute(input interface{}, op options.PluginOption) (interface{}, error) {
	domainSkip, ok := input.(*types.DomainSkip)
	if !ok {
		return nil, nil
	}
	if domainSkip.Skip {
		return nil, nil
	}
	if len(domainSkip.IP) > 1 {
		domainSkip.Skip = true
		return nil, nil
	}
	if len(domainSkip.IP) == 1 {
		key := op.TaskId + ":skipportscan:" + domainSkip.IP[0]
		flag := results.Duplicate.DuplicateLocalCache(key)
		if flag {
			rFlag := results.Duplicate.DuplicateRedisCache("TaskInfo:tmp:skipportscan:"+op.TaskId, domainSkip.IP[0])
			if rFlag {
				domainSkip.Skip = true
			}
		}
	}
	return nil, nil
}
