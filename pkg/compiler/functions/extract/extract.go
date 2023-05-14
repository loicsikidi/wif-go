package extract

import (
	"regexp"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/loicsikidi/wif-go/pkg/compiler/functions"
)

// # Extract
//
// Returns a string from a source based on an extraction template,
// which specifies the part of the attribute to extract.
//
// <string>.extract(<string>) -> <string>
//
// Examples:
//
//	'id/123456789'.extract('id/{end}')      // returns '123456789'
//	'id/123456789'.extract('{start}/')      // returns 'id'
//	'id/123456789'.extract('{all}')         // returns 'id/123456789'
//	'id/123456789'.extract('foo/{nothing}') // returns ''
//

func init() {
	functions.Register("extract", &function{})
}

type function struct{}

type Helper struct {
	Template   string
	Identifier string
	HasPrefix  bool
	HasSuffix  bool
}

func (h *Helper) GetPrefix() string {
	if h.HasPrefix && strings.Contains(h.Template, `{`) {
		return h.Template[:strings.Index(h.Template, `{`)] //nolint:gocritic
	}
	return ""
}

func (h *Helper) GetSuffix() string {
	if h.HasSuffix {
		return h.Template[strings.Index(h.Template, `}`)+1:]
	}
	return ""
}

func Identifier(source string) string {
	re := regexp.MustCompile(`{[\w-]+}`)
	return re.FindString(source)
}

func extract(source string, template string) (string, error) { //nolint:unparam
	identifier := Identifier(template)

	if identifier == "" {
		return "", nil
	}

	config := &Helper{
		Template:   template,
		Identifier: identifier,
		HasPrefix:  !strings.HasPrefix(template, identifier),
		HasSuffix:  !strings.HasSuffix(template, identifier),
	}

	if !config.HasPrefix && !config.HasSuffix {
		return source, nil
	}

	re := strings.ReplaceAll(template, identifier, "(.*)")

	r, err := regexp.Compile(re)
	if err != nil {
		return "", nil
	}
	match := r.FindStringSubmatch(source)
	if len(match) < 2 {
		return "", nil
	}

	if config.HasSuffix {
		index := strings.Index(match[1], config.GetSuffix())
		if index > 0 {
			return match[1][:index], nil
		}
	}
	return match[1], nil
}

func (f *function) GetFn() cel.EnvOption {
	return cel.Function("extract",
		cel.MemberOverload("string_extract_string",
			[]*cel.Type{cel.StringType, cel.StringType},
			cel.StringType,
			cel.BinaryBinding(func(str, extractor ref.Val) ref.Val {
				s := str.(types.String)
				ext := extractor.(types.String)
				return functions.StringOrError(extract(string(s), string(ext)))
			},
			)))
}
