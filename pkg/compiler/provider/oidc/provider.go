package oidc

import (
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/cel-go/cel"
	"github.com/loicsikidi/wif-go/pkg/common/util"
	"github.com/loicsikidi/wif-go/pkg/compiler/provider"
	pb "github.com/loicsikidi/wif-go/pkg/generated/protobuf"
	"google.golang.org/protobuf/types/known/structpb"
)

func init() {
	provider.Register("oidc", &Provider{})
}

type Provider struct{}

func (p *Provider) parseFromNativeRepresentation(raw string) (string, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(raw, jwt.MapClaims{})

	if err != nil {
		return "", err
	}

	claims, err := util.JSONEncode(token.Claims)

	if err != nil {
		return "", fmt.Errorf("error getting jwt claims")
	}
	return claims, nil
}

func (p *Provider) GetOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable("assertion", cel.DynType),
	}
}

func (p *Provider) GetInputVar(raw string) (map[string]any, error) {
	var input any

	jwtBody, err := p.parseFromNativeRepresentation(raw)

	if err == nil && jwtBody != "" {
		raw = jwtBody
	}

	if err := json.Unmarshal([]byte(raw), &input); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	_assertion, err := structpb.NewValue(input)

	if err != nil {
		return nil, fmt.Errorf("error converting given input to protobuf: %w", err)
	}

	assertion := &pb.OidcSchema{
		Assertion: _assertion,
	}

	return map[string]any{"assertion": assertion.Assertion}, nil
}
