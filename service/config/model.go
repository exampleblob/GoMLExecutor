
package config

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/viant/afs/file"
	"github.com/viant/mly/shared"
	"github.com/viant/tapper/config"
)

// Model represents model config
type Model struct {
	ID    string
	Debug bool

	Location string `json:",omitempty" yaml:",omitempty"`
	Dir      string
	URL      string

	Tags []string

	// IO and caching
	UseDict          *bool  `json:",omitempty" yaml:",omitempty"`
	DictURL          string // Deprecated: we usually extract the dictionary/vocabulary from TF graph
	shared.MetaInput `json:",omitempty" yaml:",inline"`
	OutputType       string `json:",omitempty" yaml:",omitempty"` // Deprecated: we can infer output types from TF graph
	Transformer      string `json:",omitempty" yaml:",omitempty"`

	// caching
	DataStore string `json:",omitempty" yaml:",omitempty"`

	// logging
	Stream *config.Stream `json:",omitempty" yaml:",omitempty"`

	// for health and monitoring
	Modified *Modified `json:",omitempty" yaml:",omitempty"`
	DictMeta DictionaryMeta

	Test TestPayload `json:",omitempty" yaml:",omitempty"`
}