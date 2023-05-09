package functions

import (
	"fmt"
	"sync"

	"github.com/google/cel-go/cel"
)

var (
	m         sync.Mutex
	functions = make(map[string]Interface)
)

// Interface is what function need to implement to share custom function.
type Interface interface {
	GetFn() cel.EnvOption
}

// Register is used by functions to participate in furnishing OIDC tokens.
func Register(name string, fn Interface) {
	m.Lock()
	defer m.Unlock()

	if prev, ok := functions[name]; ok {
		panic(fmt.Sprintf("duplicate function for name %q, %T and %T", name, prev, fn))
	}
	functions[name] = fn
}

// ProvideAll fetches all available functions
func ProvideAll() map[string]Interface {
	m.Lock()
	defer m.Unlock()

	return functions
}

// ProvideFrom fetches the specified function
func ProvideFrom(function string) (Interface, error) {
	m.Lock()
	defer m.Unlock()

	fn, ok := functions[function]
	if !ok {
		return nil, fmt.Errorf("%s is not a valid function", function)
	}
	return fn, nil
}
