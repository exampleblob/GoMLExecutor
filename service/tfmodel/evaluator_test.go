package tfmodel_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/viant/mly/service/tfmodel"
)

func TestBasic(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "../..")
	t.Logf("Root %s", root)
	modelDest := filepath.Join(root, "example/model/string_lookups_int_model")

	model, err := tf.LoadSavedModel(modelDest, []string{"serve"}, nil)
	assert.Nil(t, err)

	signature, err := tfmodel.Signature(model)
	ass