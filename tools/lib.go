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
		fmt.Fprintf(writer, "layer: %v hash: %v\n", l.Name, l.Hash)
	}
}

func LoadModel(ctx context.Context, URL string) (*tf.SavedModel, error) {
	fs := afs.New()

	location := url.Path(URL)
	if url.Scheme(URL, file.Scheme) != file.Scheme {
		_, name := path.Split(URL)
		location = path.Join(os.TempDir(), name)
		log.Printf("copy model files to %s", location)
		if err := fs.Copy(ctx, URL, location); err != nil {
			return nil, err
		}
	}

	model, err := tf.LoadSavedModel(location, []string{"serve"}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load model %v, due to %w", location, err)
	}
	return model, nil
}

func DiscoverDictHash(model *tf.SavedModel, writer io.Writer) error {
	signature, err := tfmodel.Signature(model)
	if err != nil {
		return err
	}

	dict, err := layers.Dictionary(model.Session, model.Graph, signature)
	if err != nil {
		return err
	}

	printDictHash(*dict, writer)

	return nil
}

func DiscoverSignature(writer io.Writer, signature *domain.Signature) error {
	encoder := sjson.NewEncoder(writer)
	return encoder.Encode(signature)
}

func DiscoverConfig(sourceURL string, model *tf.SavedModel, writer io.Writer) error {
	signature, err := tfmodel.Signature(model)
	if err != nil {
		return err
	}

	_, ID := path.Split(sourceURL)

	cfg := buildDefaultConfig(sourceURL, model, ID, signature)

	encoder := yaml.NewEncoder(writer)
	return encoder.Encode(cfg)
}

func buildDefaultConfig(sourceURL string, model *tf.SavedModel, ID string, signature *domain.Signature) endpoint.Config {
	cfg := endpoint.Config{}
	cfg.Connections = append(cfg.Connections, &datastore.Connection{
		ID:        "l1",
		Hostnames: "127.0.0.1",
		Port:      3000,
	})
	cfg.Datastores = append(cfg.Datastores, &dconfig.Datastore{
		ID:    "mly_l1",
		Cache: &scache.Config{SizeMb: 128},
		Reference: &datastore.Reference{
			Connection:   "l1",
			Namespace:    "test",
			Dataset:      ID,
			TimeToLiveMs: 0,
			RetryTimeMs:  0,
			ReadOnly:     false,
		},
		L2:       nil,
		Storable: "",
		Fields:   nil,
		Disabled: false,
	})

	cfg.Endpoint.Port = 8087
	configModel := &sconfig.Model{}
	configModel.ID = 