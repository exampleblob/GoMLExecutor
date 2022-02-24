package tfmodel

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/viant/mly/service/domain"
)

// Signature sea