package plugin

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/options"
	"github.com/Autumn-27/ScopeSentry-Scan/pkg/utils"
)

func GetName() string {
	return "demo"
}

func Install() error {
	return nil
}

func Check() error {
	return nil
}

func Uninstall() {

}

func Execute(op options.PluginOption) (interface{}, error) {
	fmt.Printf("%v", op)
	var result chan string
	go utils.Tools.ExecuteCommandToChan("whoami", []string{}, result)
	for i := range result {
		fmt.Println(i)
	}
	return nil, nil
}
