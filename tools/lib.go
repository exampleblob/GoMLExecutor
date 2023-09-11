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
	"github.com/viant/mly