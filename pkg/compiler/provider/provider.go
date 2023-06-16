package provider

import (
	"fmt"
	"sync"

	"github.com/google/cel-go/cel"
)

var (
	m         sync.Mutex
	providers = make(map[string]Provider)
)

// Provider handles CEL logic per Workload Identity Federation provider type
type Provider interface {
	GetInputVar(raw string) (map[string]any, error)
	GetOptions() []cel.EnvOption
}

func Register(name string, p Provider) {
	m.Lock()
	defer m.Unlock()

	if prev, ok := providers[name]; ok {
		panic(fmt.Sprintf("duplicate provider for name %q, %T and %T", name, prev, p))
	}
	providers[name] = p
}

func AmbientProvider(raw string) (Provider, error) {
	m.Lock()
	defer m.Unlock()

	for _, p := range providers {
		if _, err := p.GetInputVar(raw); err == nil {
			return p, nil
		}
	}

	return nil, fmt.Errorf("no provider found for given input")
}
