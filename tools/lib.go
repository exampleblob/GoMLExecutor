package tools

import (
	"context"
	"encoding/json"
	sjson "encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"reflect"
	"sort"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	sconfig "github.com/viant/mly/service/config"
	"github.com/viant/mly/service/domain"
	"github.com/viant/mly/service/endpoint"
	"github.com/viant/mly/service/layers"
	"github.com/viant/mly/service/tfmodel"
	"github.com/viant/mly/shared"
	"github.com/viant/mly/shared/common"
	dconfig "github.com/viant/mly/shared/config"
	"github.com/viant/mly/shared/config/datastore"
	"github.com/viant/scache"
	"gopkg.in/yaml.v3"
)

// Deprecated: FetchDictHash is used to inspect dictionary data from a meta URL, which is
// no longer a supported feature.
func FetchDictHash(writer io.Writer, sourceURL string, fs afs.Service) error {
	source, err := fs.DownloadWithURL(context.Background(), sourceURL)
	if err != nil {
		return err
	}
	dict := common.Dictionary{}
	if err = json.Unmarshal(source, &dict); err != nil {
		return err
	}

	printDictHash(dict, writer)

	return nil
}

func printDictHash(dict common.Dictionary, writer io.Writer) {
	fmt.Fprintf(writer, "dict hash: %v\n", dict.UpdateHash(0))
	for _, l := range dict.Layers {
		fmt.Fprintf(writ