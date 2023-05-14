package all

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
)

func getAllFunctions() []cel.EnvOption {
	opts := []cel.EnvOption{}
	for _, fn := range ProvideAll() {
		opts = append(opts, fn.GetFn())
	}
	return opts
}

var stringTests = []struct {
	expr      string
	err       string
	parseOnly bool
}{
	// Extract tests.
	{expr: `"prefix/1234567890".extract('/{id}') == "1234567890"`},
	{expr: `"prefix/1234567890/suffix".extract('/{id}') == "1234567890/suffix"`},
	{expr: `"prefix/1234567890/suffix".extract('/{id}/') == "1234567890"`},
	// Split tests.
	{expr: `"hello world".split(" ") == ["hello", "world"]`},
	{expr: `"hello world events!".split(" ", 0) == []`},
	{expr: `"hello world events!".split(" ", 1) == ["hello world events!"]`},
	{expr: `"o©o©o©o".split("©", -1) == ["o", "o", "o", "o"]`},
	// Join tests.
	{expr: `['x', 'y'].join() == 'xy'`},
	{expr: `['x', 'y'].join('-') == 'x-y'`},
	{expr: `[].join() == ''`},
	{expr: `[].join('-') == ''`},
	// Valid parse-only expressions which should generate runtime errors.
	{
		expr:      `42.extract("2") == "4"`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `"42".extract(2) == "4"`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `42.split("2") == ["4"]`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `42.split("") == ["4", "2"]`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `"42".split(2) == ["4"]`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `42.split("2", "1") == ["4"]`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `"42".split(2, 1) == ["4"]`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `"42".split("2", "1") == ["4"]`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `"42".split("2", 1, 1) == ["4"]`,
		err:       "no such overload",
		parseOnly: true,
	},
}

func TestAll(t *testing.T) {
	env, err := cel.NewEnv(getAllFunctions()...)
	if err != nil {
		t.Fatalf("cel.NewEnv(getAllFunctions()...) failed: %v", err)
	}
	for i, tst := range stringTests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			var asts []*cel.Ast
			pAst, iss := env.Parse(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("env.Parse(%v) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, pAst)
			if !tc.parseOnly {
				cAst, iss := env.Check(pAst)
				if iss.Err() != nil {
					t.Fatalf("env.Check(%v) failed: %v", tc.expr, iss.Err())
				}
				asts = append(asts, cAst)
			}
			for _, ast := range asts {
				prg, err := env.Program(ast)
				if err != nil {
					t.Fatal(err)
				}
				out, _, err := prg.Eval(cel.NoVars())
				if tc.err != "" { // nolint: gocritic
					if err == nil {
						t.Fatalf("got value %v, wanted error %s for expr: %s",
							out.Value(), tc.err, tc.expr)
					}
					if !strings.Contains(err.Error(), tc.err) {
						t.Errorf("got error %v, wanted error %s for expr: %s", err, tc.err, tc.expr)
					}
				} else if err != nil {
					t.Fatal(err)
				} else if out.Value() != true {
					t.Errorf("got %v, wanted true for expr: %s", out.Value(), tc.expr)
				}
			}
		})
	}
}
