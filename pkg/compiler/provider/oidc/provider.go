package oidc

import (
	"encoding/json"
	"fmt"

	"github.com/google/cel-go/cel"
	pb "github.com/loicsikidi/wif-go/pkg/generated/protobuf"
	"google.golang.org/protobuf/types/known/structpb"
)

type Provider struct{}

func (p *Provider) GetOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable("assertion", cel.DynType),
	}
}

func (p *Provider) GetInputVar(raw string) (map[string]any, error) {
	var _json any

	if err := json.Unmarshal([]byte(raw), &_json); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	_assertion, err := structpb.NewValue(_json)

	if err != nil {
		return nil, fmt.Errorf("error converting JSON to protobuf: %w", err)
	}

	assertion := &pb.OidcSchema{
		Assertion: _assertion,
	}

	return map[string]any{"assertion": assertion.Assertion}, nil
}
