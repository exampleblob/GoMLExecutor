package service

import (
	"compress/gzip"
	"context"
	sjson "encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"runtime/trace"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	tf "github.com/tensorflow/