package saml

import (
	"encoding/xml"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/loicsikidi/wif-go/pkg/compiler/provider"
	pb "github.com/loicsikidi/wif-go/pkg/generated/protobuf"
	"github.com/loicsikidi/wif-go/pkg/saml"
	"google.golang.org/protobuf/types/known/structpb"
)

func init() {
	provider.Register("saml", &Provider{})
}

const (
	assertion string = "assertion"
)

type Provider struct{}

func (p *Provider) GetOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable(assertion, cel.ObjectType("wifgo.v1beta.SamlSchema")),
		cel.Types(&pb.SamlSchema{}),
	}
}

func (p *Provider) GetInputVar(raw string) (map[string]any, error) {
	var input saml.Response

	err := xml.Unmarshal([]byte(raw), &input)

	if err != nil {
		return nil, fmt.Errorf("error unable to parse SAML response: %w", err)
	}

	attrMap := input.GetMapAttributes()
	attributes, err := structpb.NewValue(attrMap[saml.AttributesStr])

	if err != nil {
		return nil, fmt.Errorf("error converting given input to protobuf: %w", err)
	}

	return map[string]any{
		assertion: &pb.SamlSchema{
			Subject:    attrMap[saml.SubjectStr].(string),
			Attributes: attributes,
		},
	}, nil
}
