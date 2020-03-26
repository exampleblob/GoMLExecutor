package main

import (
	"os"

	"github.com/viant/mly/example/client"
	slfmodel "github.com/viant/mly/example/transformer/slf/model"
	"github.com/viant/mly/shared/common"
	"github.com/viant/mly/shared/common/storable"
)

func main() {
	storableSrv := storable.Singleton()

	// in actual client code, any type information should be available within the context of the caller, so 