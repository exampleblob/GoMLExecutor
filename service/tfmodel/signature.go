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
		return nil, fmt.Errorf("failed to lookup signature: %v", domain.DefaultSignatureKey)
	}

	sig := &domain.Signature{
		Method: signature.MethodName,
	}

	for layerName, tfInfo := range signature.Outputs {
		output := domain.Output{
			Name: layerName,
		}

		// tfInfo.Name is the Tensor name
		operationName := tfInfo.Name
		if index := strings.Index(operationName, ":"); index != -1 {
			indexValue := operationName[index+1:]
			operationName = operationName[:index]
			output.Index, _ = strconv.Atoi(indexValue)
		}

		if output.Operation = model.Graph.Operation(operationName); output.Oper