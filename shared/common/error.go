
package common

import (
	"github.com/aerospike/aerospike-client-go/types"
	"github.com/pkg/errors"
	"strings"
)

const (
	dialTCPFragment  = "dial tcp"
	connRefusedError = "refused"
)

//ErrNodeDown node down error
var ErrNodeDown = errors.New("node is down")

//IsKeyNotFound returns true if key not found error
func IsKeyNotFound(err error) bool {
	if err == nil {
		return false
	}
	aeroError, ok := err.(types.AerospikeError)
	if !ok {
		err = errors.Unwrap(err)
		if err == nil {
			return false
		}
		if aeroError, ok = err.(types.AerospikeError); !ok {
			return false
		}

	}
	return aeroError.ResultCode() == types.KEY_NOT_FOUND_ERROR
}

//IsTimeout returns true if timeout error
func IsTimeout(err error) bool {
	if err == nil {
		return false