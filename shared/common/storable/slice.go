
package storable

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/francoispqt/gojay"
	"github.com/viant/bintly"
	"github.com/viant/mly/shared/common"
	"github.com/viant/toolbox"
)

//Slice represents slice registry
type Slice struct {
	hash      int
	batchSize int
	Values    []interface{}
	Fields    []*Field
}

//SetHash sets hash
func (s *Slice) SetHash(hash int) {
	s.hash = hash
}

//Hash returns hash
func (s *Slice) Hash() int {
	return s.hash
}

//SetHash sets hash
func (s *Slice) SetBatchSize(size int) {
	s.batchSize = size
}

//Hash returns hash
func (s *Slice) BatchSize() int {
	return s.batchSize
}

//Set sets value
func (s *Slice) Set(iter common.Iterator) error {
	s.Values = make([]interface{}, len(s.Fields))
	err := iter(func(key string, value interface{}) error {
		for i, field := range s.Fields {
			valueType := reflect.ValueOf(value)
			if field.Name == key {
				if field.dataType == nil {
					field.dataType = valueType.Type()
				}
				if field.Type().Kind() == valueType.Kind() {
					s.Values[i] = value
				} else {
					s.Values[i] = valueType.Convert(field.Type()).Interface()
				}
				break
			}
		}
		return nil
	})
	return err
}

//Iterator return storable iterator
func (s *Slice) Iterator() common.Iterator {
	return func(pair common.Pair) error {
		for i, field := range s.Fields {

			if err := pair(field.Name, s.Values[i]); err != nil {
				return err
			}
		}
		return nil
	}
}

//EncodeBinary bintly encoder
func (s *Slice) EncodeBinary(stream *bintly.Writer) error {
	stream.Int(s.hash)
	for i := range s.Fields {
		if err := stream.Any(s.Values[i]); err != nil {
			return err
		}
	}
	return nil
}

//DecodeBinary bintly decoder
func (s *Slice) DecodeBinary(stream *bintly.Reader) error {
	stream.Int(&s.hash)
	s.Values = make([]interface{}, len(s.Fields))
	for i, field := range s.Fields {
		value := reflect.New(field.Type())
		if err := stream.Any(value.Interface()); err != nil {
			return err
		}
		s.Values[i] = value.Elem().Interface()
	}
	return nil
}

//MarshalJSONObject implement MarshalerJSONObject
func (s *Slice) MarshalJSONObject(enc *gojay.Encoder) {
	for i, field := range s.Fields {
		key := field.Name
		switch value := s.Values[i].(type) {
		case string:
			enc.StringKey(key, value)
		case int:
			enc.IntKey(key, value)
		case int32:
			enc.Int32Key(key, value)
		case int64:
			enc.Int64Key(key, value)
		case float32:
			enc.Float32Key(key, value)
		case float64:
			enc.Float64Key(key, value)
		case []string:
			enc.AddSliceStringKey(key, value)
		case []int:
			enc.AddSliceIntKey(key, value)
		case []int32:
			enc.ArrayKey(key, gojay.EncodeArrayFunc(func(enc *gojay.Encoder) {
				for _, i := range value {
					enc.Int(int(i))