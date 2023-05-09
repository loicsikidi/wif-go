package provider

import "github.com/google/cel-go/cel"

const (
	OIDC = iota
)

type Backend = int

// Provider handles CEL logic per Workload Identity Federation provider type
type Provider interface {
	GetOptions() []cel.EnvOption
	GetInputVar(raw string) (map[string]any, error)
}
