package tfmodel

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/viant/mly/service/domain"
)

// Signature searches the Tensorflow operation graph for inputs and outputs.
func Signature(model *tf.SavedModel) (*domain.Signature, error) {
	signature, ok := model.Signatures[domain.DefaultSignatureKey]
	if !ok {
