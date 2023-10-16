package tools

import "fmt"

type Options struct {
	Mode      string `short:"m" long:"mode" choice:"discover" choice:"run" description:"mode"`
	