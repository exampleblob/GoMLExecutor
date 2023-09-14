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

// Deprecated: FetchDictHash 