package client

import (
	"log"

	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/common"
)

type fieldOffset int

const (
	// oov = out of vocabulary
	oovString = "[UNK]"
	oovInt    = 0

	defaultPrec = 10

	unknownKeyField = fiel