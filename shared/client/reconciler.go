package client

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"

	"github.com/viant/xunsafe"
)

// ReconcileData reconciles target with cached and predicted data
// target is either the pointer to the result or a pointer to a slice of results from
// the prediction server
func reconcileData(prefix string, target interface{}, cachable Cachable, cached []interface{}) error {
	targetType := reflect.TypeOf(target).Elem()
	targetPtr := xunsafe.AsPointer(target)

	if prefix != "" {
		log.Printf("%s reconciling: %T %+v", p