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
		log.Printf("%s reconciling: %T %+v", prefix, target, target)
	}

	switch targetType.Kind() {
	case reflect.Struct:
		if !cachable.CacheHit(0) {
			// the target memory already has actual value
			return nil
		}

		// directly replace the target memory with cached value
		*(*unsafe.Pointer)(targetPtr) = *(*unsafe.Pointer)(xunsafe.AsPointer(cached[0]))
		return nil
	case reflect.Slice:
		// noop
	default:
		return fmt.Errorf("unsupported target type