
package tools

import (
	"context"
	"fmt"

	"github.com/jessevdk/go-flags"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/viant/afs"

	slog "log"
	"strings"

	"github.com/viant/mly/service/endpoint"
	"github.com/viant/mly/service/layers"
	"github.com/viant/mly/service/tfmodel"
	"github.com/viant/mly/shared/common"
	"github.com/viant/tapper/config"
	"github.com/viant/tapper/io"
	"github.com/viant/tapper/log"
	"github.com/viant/tapper/msg"
	"github.com/viant/tapper/msg/json"
)

func Run(args []string) {
	options := &Options{}

	_, err := flags.ParseArgs(options, args)
	if err != nil {
		slog.Fatal(err)
	}
	if err := options.Validate(); err != nil {
		slog.Fatal(err)
	}

	switch options.Mode {
	case "discover":
		err = Discover(options)
		if err != nil {
			slog.Fatal(err)
		}
	case "run":
		cfg, err := endpoint.NewConfigFromURL(context.Background(), options.ConfigURL)
		if err != nil {
			slog.Fatal(err)
			return
		}
		srv, err := endpoint.New(cfg)
		if err != nil {
			slog.Fatal(err)
			return
		}

		srv.ListenAndServe()
	}
}

func Discover(options *Options) error {
	fs := afs.New()
	writer, err := GetWriter(options.DestURL, fs)
	if err != nil {
		return err
	}