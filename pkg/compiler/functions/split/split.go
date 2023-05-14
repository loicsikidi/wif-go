package split

import (
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/loicsikidi/wif-go/pkg/compiler/functions"
)

// # Split
//
// Returns a list of strings split from the input by the given separator. The function accepts
// an optional argument specifying a limit on the number of substrings produced by the split.
//
// When the split limit is 0, the result is an empty list. When the limit is 1, the result is the
// target string to split. When the limit is a negative number, the function behaves the same as
// split all.
//
//	<string>.split(<string>) -> <list<string>>
//	<string>.split(<string>, <int>) -> <list<string>>
//
// Examples:
//
//	'hello hello hello'.split(' ')     // returns ['hello', 'hello', 'hello']
//	'hello hello hello'.split(' ', 0)  // returns []
//	'hello hello hello'.split(' ', 1)  // returns ['hello hello hello']
//	'hello hello hello'.split(' ', 2)  // returns ['hello', 'hello hello']
//	'hello hello hello'.split(' ', -1) // returns ['hello', 'hello', 'hello']
//

func init() {
	functions.Register("split", &function{})
}

type function struct{}

func split(str, sep string) ([]string, error) { //nolint:unparam
	return strings.Split(str, sep), nil
}

func splitN(str, sep string, n int64) ([]string, error) { //nolint:unparam
	return strings.SplitN(str, sep, int(n)), nil
}

func (f *function) GetFn() cel.EnvOption {
	return cel.Function("split",
		cel.MemberOverload("string_split_string", []*cel.Type{cel.StringType, cel.StringType}, cel.ListType(cel.StringType),
			cel.BinaryBinding(func(str, separator ref.Val) ref.Val {
				s := str.(types.String)
				sep := separator.(types.String)
				return functions.ListStringOrError(split(string(s), string(sep)))
			})),
		cel.MemberOverload("string_split_string_int", []*cel.Type{cel.StringType, cel.StringType, cel.IntType}, cel.ListType(cel.StringType),
			cel.FunctionBinding(func(args ...ref.Val) ref.Val {
				s := args[0].(types.String)
				sep := args[1].(types.String)
				n := args[2].(types.Int)
				return functions.ListStringOrError(splitN(string(s), string(sep), int64(n)))
			})))
}
