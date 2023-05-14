package split

import (
	"fmt"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/loicsikidi/wif-go/pkg/compiler/functions"
)

// # Join
//
// Returns a new string where the elements of string list are concatenated.
//
// The function also accepts an optional separator which is placed between elements in the resulting string.
//
// <list<string>>.join() -> <string>
// <list<string>>.join(<string>) -> <string>
//
// Examples:
//
//	['hello', 'mellow'].join() // returns 'hellomellow'
//	['hello', 'mellow'].join(' ') // returns 'hello mellow'
//	[].join() // returns ''
//	[].join('/') // returns ''
//

func init() {
	functions.Register("join", &function{})
}

type function struct{}

func joinValSeparator(strs traits.Lister, separator string) (string, error) {
	sz := strs.Size().(types.Int)
	var sb strings.Builder
	for i := types.Int(0); i < sz; i++ {
		if i != 0 {
			sb.WriteString(separator)
		}
		elem := strs.Get(i)
		str, ok := elem.(types.String)
		if !ok {
			return "", fmt.Errorf("join: invalid input: %v", elem)
		}
		sb.WriteString(string(str))
	}
	return sb.String(), nil
}

func (f *function) GetFn() cel.EnvOption {
	return cel.Function("join",
		cel.MemberOverload("list_join", []*cel.Type{cel.ListType(cel.StringType)}, cel.StringType,
			cel.UnaryBinding(func(list ref.Val) ref.Val {
				l := list.(traits.Lister)
				return functions.StringOrError(joinValSeparator(l, ""))
			})),
		cel.MemberOverload("list_join_string", []*cel.Type{cel.ListType(cel.StringType), cel.StringType}, cel.StringType,
			cel.BinaryBinding(func(list, delim ref.Val) ref.Val {
				l := list.(traits.Lister)
				d := delim.(types.String)
				return functions.StringOrError(joinValSeparator(l, string(d)))
			})))
}
