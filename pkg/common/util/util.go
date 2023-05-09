package util

import (
	"encoding/json"
	"errors"
	"reflect"
)

// Any returns true if one of the strings in the slice satisfies the predicate f.
func Any(vs []string, f func(any) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

// Map returns a new slice containing the results of applying the function f to each string in the original slice.
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

// Filter returns a new slice containing all strings in the slice that satisfy the predicate f.
func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// JSONEncode encodes a value to JSON
func JSONEncode(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ConvertStringMapToAny converts a map of string elements to a map of any elements
func ConvertStringMapToAny(m map[string]string) map[string]any {
	newMap := make(map[string]any)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// MergeMaps merges two maps
func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	merged := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

// GetMapKeys returns the keys of a map
func GetMapKeys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// GetMapKeys returns the keys of a map
func GetMapValues[K comparable, V any](m map[K]V) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

var (
	ErrSourceNotArray = errors.New("source value is not an array")
	ErrReducerNil     = errors.New("reducer function cannot be nil")
	ErrReducerNotFunc = errors.New("reducer argument must be a function")
)

// Reduce an array of something into another thing
func Reduce(source, initialValue, reducer any) (any, error) {
	srcV := reflect.ValueOf(source)
	kind := srcV.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, ErrSourceNotArray
	}
	if reducer == nil {
		return nil, ErrReducerNil
	}
	rv := reflect.ValueOf(reducer)
	if rv.Kind() != reflect.Func {
		return nil, ErrReducerNotFunc
	}
	// copy initial value as accumulator, and get the reflection value
	accumulator := initialValue
	accV := reflect.ValueOf(accumulator)
	for i := 0; i < srcV.Len(); i++ {
		entry := srcV.Index(i)
		// call reducer via reflection
		reduceResults := rv.Call([]reflect.Value{
			accV,               // send accumulator value
			entry,              // send current source entry
			reflect.ValueOf(i), // send current loop index
		})
		accV = reduceResults[0]
	}
	return accV.Interface(), nil
}
