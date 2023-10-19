package tools

import "fmt"

type Options struct {
	Mode      string `short:"m" long:"mode" choice:"discover" choice:"run" description:"mode"`
	Operation string `short:"o" long:"opt" choice:"dictHash" choice:"signature" choice:"layers" choice:"config"`
	SourceURL string `short:"s" long:"src" description:"source location, required by all opt where mode=discover"`
	DestURL   string `short:"d" long:"dest" description:"dest location, required if opt=layers"`
	ConfigURL string `short:"c" long:"config" description:"required if mode=run"`
}

fun