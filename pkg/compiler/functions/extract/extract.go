package extract

import (
	"regexp"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/loicsikidi/wif-go/pkg/compiler/functions"
)

func init() {
	functions.Register("extract", &extract{})
}

type extract struct{}

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

func Extract(source string, template string) string {
	identifier := Identifier(template)

	if identifier == "" {
		return ""
	}

	config := &Helper{
		Template:   template,
		Identifier: identifier,
		HasPrefix:  !strings.HasPrefix(template, identifier),
		HasSuffix:  !strings.HasSuffix(template, identifier),
	}

	if !config.HasPrefix && !config.HasSuffix {
		return source
	}

	re := strings.ReplaceAll(template, identifier, "(.*)")

	r, err := regexp.Compile(re)
	if err != nil {
		return ""
	}
	match := r.FindStringSubmatch(source)
	if len(match) < 2 {
		return ""
	}

	if config.HasSuffix {
		index := strings.Index(match[1], config.GetSuffix())
		if index > 0 {
			return match[1][:index]
		}
	}
	return match[1]
}

func (f *extract) GetFn() cel.EnvOption {
	return cel.Function("extract",
		cel.MemberOverload("string_extract_string",
			[]*cel.Type{cel.StringType, cel.StringType},
			cel.StringType,
			cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return types.String(
					Extract(lhs.Value().(string), rhs.Value().(string)))
			},
			)))
}
