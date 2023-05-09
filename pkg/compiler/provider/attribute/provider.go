package attribute

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/loicsikidi/wif-go/pkg/common/util"
	pb "github.com/loicsikidi/wif-go/pkg/generated/protobuf"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	Attribute string = "attribute"
	Google    string = "google"
	Assertion string = "assertion"
)

var (
	VARIABLES = []string{Attribute, Google, Assertion}
)

// GetAttributeName is an helper function to get a derived attribute
func GetAttributeName(attr string) string {
	return fmt.Sprintf("%s.%s", Attribute, attr)
}

// GetAssertionName is an helper function to get an assertion
func GetAssertionName(ass string) string {
	return fmt.Sprintf("%s.%s", Assertion, ass)
}

func GetAttributeInputVar(customMap map[string]any) (string, error) {
	attrMap := map[string]map[string]any{}
	for _, v := range VARIABLES {
		attrMap[v] = map[string]any{}
	}

	for k, v := range customMap {
		isDerivedClaim := util.Any([]string{Attribute, Google}, func(s any) bool {
			return strings.HasPrefix(k, fmt.Sprintf("%s.", s))
		})
		if isDerivedClaim {
			family := strings.Split(k, ".")[0]
			name := strings.Split(k, ".")[1]
			attrMap[family][name] = v
		} else {
			attrMap[k] = v.(map[string]any)
		}
	}

	json, err := util.JSONEncode(attrMap)

	if err != nil {
		return "", fmt.Errorf("error encoding custom attribute map into JSON: %w", err)
	}

	return json, nil
}

type Provider struct{}

func (p *Provider) GetOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Variable(Attribute, cel.DynType),
		cel.Variable(Google, cel.DynType),
		cel.Variable(Assertion, cel.DynType),
	}
}

func (p *Provider) GetInputVar(raw string) (map[string]any, error) {
	var _json any

	if err := json.Unmarshal([]byte(raw), &_json); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	attribute := &pb.AttributeSchema{}
	for k, v := range _json.(map[string]any) {
		pb, err := structpb.NewValue(v)

		if err != nil {
			return nil, fmt.Errorf("error converting JSON to protobuf: %w", err)
		}

		if k == Google {
			attribute.Google = pb
		}
		if k == Assertion {
			attribute.Assertion = pb
		}
		if k == Attribute {
			attribute.Attribute = pb
		}
	}

	return map[string]any{
		Attribute: attribute.Attribute,
		Assertion: attribute.Assertion,
		Google:    attribute.Google}, nil
}
