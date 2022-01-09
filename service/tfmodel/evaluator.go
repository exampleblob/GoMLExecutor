package tfmodel

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"runtime/trace"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/viant/mly/service/clienterr"
	"githu