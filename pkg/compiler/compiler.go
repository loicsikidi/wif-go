package compiler

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/loicsikidi/wif-go/pkg/common/util"
	"github.com/loicsikidi/wif-go/pkg/compiler/provider"
	"github.com/loicsikidi/wif-go/pkg/compiler/provider/attribute"

	// Link in all of the functions
	allFns "github.com/loicsikidi/wif-go/pkg/compiler/functions/all"
	allProviders "github.com/loicsikidi/wif-go/pkg/compiler/provider/all"
)

// Standard google attributes defined by Google Cloud Platform.
const (
	// Attribute representing an unique identifier of an identity.
	GoogleSubject = "google.subject"
	// Attribute representing a set of groups that the identity belongs to.
	GoogleGroups = "google.groups"
)

// Limitations set by Google Cloud Platform.
//
// (See more at https://cloud.google.com/iam/quotas#limits)
const (
	// google.subject can't exceed 127 bytes
	MaximumSubjectLengthInBytes = 127
	// custom attributes can't exceed 100 characters
	MaximumCustomAttributeNameSize = 100
	// Custom attributes can't exceed 8192 bytes
	MaximumCustomAttributesLengthInBytes = 8192
	// Attribute mapping expression can't exceed 2048 bytes
	MaximumAttributeExpressionLengthInBytes = 2048
	// Attribute condition expression can't exceed 4096 bytes
	MaximumAttributeConditionLengthInBytes = 4096
	// Workload identity federation can't exceed 50 custom attributes
	MaximumCustomAttributes = 50
)

// ErrAttrConditionFailed means that a credential
// was rejected by the attribute condition.
var ErrAttrConditionFailed = errors.New("the given credential is rejected by the attribute condition")

// Input used to compile a Workload Identity Federation (WIF) expression
type Input struct {
	// [Required] Payload is the source of the expression, in this project it's an external token (eg. a JWT, a SAML2.0 response, etc.)
	Payload string
	// [Required] AttributeMapping defines how to derive the value from an external token into attributes interpretable by GCP IAM.
	// The map key is the target attribute (eg. google.subject, google.groups, etc.).
	// The map value is a Common Expression Language (CEL) expression that transforms one or more attributes from the external token (.
	AttributeMapping map[string]string
	// [Optional] AttributeCondition is CEL expression that can check assertion attributes and target attributes (eg. 'admins' in google.groups).
	// If the attribute condition evaluates to true for a given credential, the credential is accepted.
	// Otherwise, the credential is rejected.
	AttributeCondition string
}

// Compiler is the main struct used to compile a WIF expression
type Compiler struct {
	// Input used to compile a Workload Identity Federation expression
	Input *Input
	// Target Provider supported by Workload Identity Federation  (eg. OIDC, SAML, etc.)
	Provider provider.Provider
}

// eval is a helper function that compiles and evaluates a CEL expression
func eval(env *cel.Env, input map[string]any, expr string) (ref.Val, error) {
	if strings.Contains(expr, "timestamp(int(") {
		return nil, fmt.Errorf("create a timestamp using unix timestamp is not currently supported by the Workload Identity Federation CEL implementation")
	}

	ast, issues := env.Compile(expr)

	if issues.Err() != nil {
		return nil, fmt.Errorf("error compiling CEL expression: %w", issues.Err())
	}

	prg, _ := env.Program(ast)
	result, _, err := prg.Eval(input)

	if err != nil {
		return nil, fmt.Errorf("error evaluating CEL expression: %w", err)
	}

	return result, nil
}

// addCustomFn adds Workload Identity Federation custom functions to the CEL environment
func addCustomFn(opts []cel.EnvOption) []cel.EnvOption {
	for _, fn := range allFns.ProvideAll() {
		opts = append(opts, fn.GetFn())
	}
	return opts
}

// getCustomAttr returns a list of custom attributes names
func getCustomAttr(attributes map[string]any) []string {
	return util.Filter(util.GetMapKeys(attributes), func(v string) bool {
		return strings.HasPrefix(v, fmt.Sprintf("%s.", attribute.Attribute))
	})
}

// convertProtoMapToRegularMap converts a map of proto to
// a map of regular type (eg. string, int32, float64, etc.)
func convertProtoMapToRegularMap(input map[string]any) (map[string]any, error) {
	var _json any

	rawJSON, err := util.JSONEncode(input)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(rawJSON), &_json); err != nil {
		return nil, err
	}

	return _json.(map[string]any), nil
}

// checkGoogleGroupsValue checks if the value of the google.groups attribute is valid
func checkGoogleGroupsValue(list ref.Val) error {
	if list.Type() != types.ListType {
		return fmt.Errorf("the mapped attribute '%s' must be of type LIST<STRING>", GoogleGroups)
	}
	it := list.(traits.Lister).Iterator()
	var i = int64(0)
	for ; it.HasNext() == types.Bool(true); i++ {
		elem := it.Next()
		if elem.Type() != types.StringType {
			return fmt.Errorf("the elements in mapped attribute '%s' must be of type STRING", GoogleGroups)
		}
	}
	return nil
}

// Run compiles a Workload Identity Federation expression and returns a map of derived attributes
func (c *Compiler) Run() (map[string]any, error) {
	derivedAttributes := map[string]any{}

	// Input validation
	if c.Input == nil || c.Input.Payload == "" || c.Input.AttributeMapping == nil {
		return nil, fmt.Errorf("input is invalid. Payload and AttributeMapping are required")
	}

	// Input.AttributeMapping validation
	for k := range c.Input.AttributeMapping {
		if !strings.HasPrefix(k, fmt.Sprintf("%s.", attribute.Attribute)) && k != GoogleSubject && k != GoogleGroups {
			return nil, fmt.Errorf("invalid attribute mapping key: %s.\nOnly 'google.subject', 'google.groups' and 'attribute.<custom_attribute>' are accepted", k)
		}
	}

	if c.Provider == nil {
		provider, err := allProviders.AmbientProvider(c.Input.Payload)

		if err != nil {
			return nil, err
		}

		c.Provider = provider
	}

	input, err := c.Provider.GetInputVar(c.Input.Payload)

	if err != nil {
		return nil, err
	}

	env, err := cel.NewEnv(addCustomFn(c.Provider.GetOptions())...)

	if err != nil {
		return nil, fmt.Errorf("error creating CEL environment: %w", err)
	}

	if err := c.preValidation(); err != nil {
		return nil, err
	}

	for k, v := range c.Input.AttributeMapping {
		val, err := eval(env, input, v)

		if err != nil {
			return nil, err
		}

		if k != GoogleGroups {
			// nominal case: we expect a string
			if val.Type() != types.StringType {
				return nil, fmt.Errorf("the mapped attribute '%s' must be of type STRING", k)
			}
			output := val.(types.String)
			derivedAttributes[k] = string(output)
		} else {
			// special case: we expect a list of strings
			// The elements in mapped attribute 'google.groups' must be of type STRING.
			if err := checkGoogleGroupsValue(val); err != nil {
				return nil, err
			}

			output, err := val.ConvertToNative(reflect.TypeOf([]any{}))

			if err != nil {
				return nil, err
			}

			derivedAttributes[k] = output
		}
	}

	if err := c.postValidation(derivedAttributes); err != nil {
		return nil, err
	}

	if c.Input.AttributeCondition != "" {
		newInput, err := convertProtoMapToRegularMap(input)

		if err != nil {
			return nil, err
		}

		mergedMap := util.MergeMaps(newInput, derivedAttributes)

		attributeProvider := attribute.Provider{}
		inputVar, err := attribute.GetAttributeInputVar(mergedMap)

		if err != nil {
			return nil, fmt.Errorf("error producing attribute input var: %w", err)
		}

		attrInput, err := attributeProvider.GetInputVar(inputVar)

		if err != nil {
			return nil, err
		}

		attrEnv, err := cel.NewEnv(addCustomFn(attributeProvider.GetOptions())...)

		if err != nil {
			return nil, fmt.Errorf("error creating attribute condition CEL environment: %w", err)
		}

		val, err := eval(attrEnv, attrInput, c.Input.AttributeCondition)

		if err != nil {
			return nil, err
		}

		if val.Type() != types.BoolType {
			return nil, ErrAttrConditionFailed
		}

		condition := val.(types.Bool)
		if !bool(condition) {
			return nil, ErrAttrConditionFailed
		}
	}
	return derivedAttributes, nil
}

// preValidation validates attribute mapping's conformity
func (c *Compiler) preValidation() error {
	for _, expr := range c.Input.AttributeMapping {
		if len(expr) > MaximumAttributeExpressionLengthInBytes {
			return fmt.Errorf("the maximum length of an attribute mapping expression is %d characters", MaximumAttributeExpressionLengthInBytes)
		}
	}

	if len(c.Input.AttributeCondition) > MaximumAttributeConditionLengthInBytes {
		return fmt.Errorf("the maximum length of an attribute condition expression is %d characters", MaximumAttributeConditionLengthInBytes)
	}

	customAttr := getCustomAttr(util.ConvertStringMapToAny(c.Input.AttributeMapping))

	for _, attr := range util.Map(customAttr, func(v string) string {
		return strings.Split(v, ".")[1]
	}) {
		r := regexp.MustCompile("^[a-z0-9_]{1,100}$")
		if !r.MatchString(attr) {
			return fmt.Errorf("invalid mapped attribute key: %s. The maximum length of a mapped attribute key is 100 characters and may only contain the characters [a-z0-9_]", attr)
		}
	}

	return nil
}

// postValidation validates derived attributes's conformity
func (c *Compiler) postValidation(derivedAttributes map[string]any) error {
	if _, ok := derivedAttributes[GoogleSubject]; !ok {
		return fmt.Errorf("missing '%s' attribute", GoogleSubject)
	}

	if len(derivedAttributes[GoogleSubject].(string)) > MaximumSubjectLengthInBytes {
		return fmt.Errorf("the size of mapped attribute '%s' exceeds the %d bytes limit", GoogleSubject, MaximumSubjectLengthInBytes)
	}

	mappedAttrSize, _ := util.Reduce(util.GetMapValues(derivedAttributes), 0, func(acc int, currentValue any, currentIndex int) int {
		switch val := currentValue.(type) {
		case string:
			acc += len(val)
		case []any:
			for _, v := range val {
				acc += len(v.(string))
			}
		}
		return acc
	})

	if mappedAttrSize.(int) > MaximumCustomAttributesLengthInBytes {
		return fmt.Errorf("the size of mapped attributes exceeds the %d bytes limit", MaximumCustomAttributesLengthInBytes)
	}

	if len(getCustomAttr(derivedAttributes)) > MaximumCustomAttributes {
		return fmt.Errorf("custom attributes are limited to %d", MaximumCustomAttributes)
	}

	return nil
}
