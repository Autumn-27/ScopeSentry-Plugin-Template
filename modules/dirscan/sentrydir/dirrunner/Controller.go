// Package dirrunner-----------------------------
// @file      : runner.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2024/4/28 23:37
// -------------------------------------------
package dirrunner

import (
	"fmt"
	"github.com/Autumn-27/ScopeSentry-Scan/internal/types"
	"github.com/Autumn-27/ScopeSentry-Scan/modules/dirscan/sentrydir/dircore"
	"sort"
	"strconv"
	"strings"
)

type Controller struct {
	Targets    []string
	Dictionary string
	Request    dircore.Request
	Threads    int
}

func (c *Controller) Run(options dircore.Options) {
	if options.Extensions == nil {
		options.Extensions = []string{"php", "aspx", "jsp", "html", "js"}
	}
	if options.IncludeStatusCodes == nil {
		statusCodes := ParseStatusCodes("200-399,401,403,500-520")
		options.IncludeStatusCodes = statusCodes
	}
	if options.MatchCallback == nil {
		options.MatchCallback = func(response types.HttpResponse) {
			fmt.Printf("%v - %v\n", response.StatusCode, response.Url)
		}
	}
	for _, target := range c.Targets {
		c.SetUrl(target)
		fuzz := dircore.Fuzzer{
			Dictionary:         c.Dictionary,
			Threads:            c.Threads,
			Request:            c.Request,
			BasePath:           "",
			Options:            options,
			MaxSameLen:         30,
			ResponseCodeLength: make(map[string]int),
			Ct:                 options.Ct,
		}
		fuzz.Start()
	}
}

func (c *Controller) SetUrl(target string) {
	if !strings.HasSuffix(target, "/") {
		target += "/"
	}
	c.Request.Url = target
}

func ParseStatusCodes(statusCodesString string) []int {
	statusCodes := []int{}

	ranges := strings.Split(statusCodesString, ",")

	for _, rangeStr := range ranges {
		if strings.Contains(rangeStr, "-") {
			rangeParts := strings.Split(rangeStr, "-")
			if len(rangeParts) != 2 {
				fmt.Println("Invalid range:", rangeStr)
				continue
			}
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				fmt.Println("Invalid start value:", rangeStr)
				continue
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				fmt.Println("Invalid end value:", rangeStr)
				continue
			}
			for i := start; i <= end; i++ {
				statusCodes = append(statusCodes, i)
			}
		} else {
			code, err := strconv.Atoi(rangeStr)
			if err != nil {
				fmt.Println("Invalid code:", rangeStr)
				continue
			}
			statusCodes = append(statusCodes, code)
		}
	}
	sort.Ints(statusCodes)
	return statusCodes
}
