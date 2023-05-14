package compiler

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/loicsikidi/wif-go/pkg/compiler/provider"
	"github.com/loicsikidi/wif-go/pkg/compiler/provider/attribute"
	"github.com/loicsikidi/wif-go/pkg/compiler/provider/oidc"
)

const (
	invalidCelExpr = "invalid.cel.expr"
	jwtPayloadBody = `{"sub": "1234567890", "is_admin": true, "iat": 1683438895, "groups": ["group1", "group2"]}`
)

func generateStr(length int) string {
	s := ""
	for i := 0; i < length; i++ {
		s += "a"
	}
	return s
}

func generateAttrMap(length int) map[string]string {
	m := map[string]string{GoogleSubject: `'1234567890'`}
	for i := 0; i < length; i++ {
		m[attribute.GetAttributeName(strconv.Itoa(i))] = `'Hello World'`
	}
	return m
}

func TestInvalidInputs(t *testing.T) {
	tests := []struct {
		input    *Input
		provider provider.Provider
	}{
		{
			input: &Input{
				Payload: `{"sub": "prefix/1234567890", "name": "John Doe", "iat": 1516239022}`,
			},
			provider: &oidc.Provider{},
		},
		{
			input: &Input{
				AttributeMapping: map[string]string{GoogleSubject: "assertion.sub.extract('/{id}')"},
			},
			provider: &oidc.Provider{},
		},
		{
			input:    &Input{},
			provider: &oidc.Provider{},
		},
	}

	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			c := &Compiler{
				Input:    tc.input,
				Provider: tc.provider,
			}
			_, err := c.Run()

			if err == nil {
				t.Fatalf("Run(%v) -> expect exception", err)
			}
		})
	}
}

func TestAttributeMapping(t *testing.T) {
	type Expected struct {
		IsError bool
	}
	tests := []struct {
		input    *Input
		provider provider.Provider
		expected *Expected
	}{
		// Test success
		{
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
					GoogleGroups:  attribute.GetAssertionName("groups"),
				},
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: false,
			},
		},
		// Test failure when the mapped value is not a string
		{
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleSubject: "true",
				},
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
		// Test failure when google.groups attr is not a list
		{
			input: &Input{
				Payload: `{"sub":"1234567890", "groups": "group1"}`,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
					GoogleGroups:  attribute.GetAssertionName("groups"),
				},
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
		// Test failure when google.groups attr is not a list of strings
		{
			input: &Input{
				Payload: `{"sub":"1234567890", "groups": ["group1", "group2", 1]}`,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
					GoogleGroups:  attribute.GetAssertionName("groups"),
				},
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
		// Test failure when the payload is an invalid JSON
		{
			input: &Input{
				Payload: `{sub: "1234567890"}`,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
				},
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
		// Test failure when the mapped key is unauthorized
		{
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					"invalid_key": attribute.GetAssertionName("sub"),
				},
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
		// Test failure when missing subject attribute
		{
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleGroups: attribute.GetAssertionName("groups"),
				},
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
		// Test failure when the attribute subject value is too big
		{
			input: &Input{
				Payload: fmt.Sprintf(`{"sub": "%s"}`, generateStr(MaximumSubjectLengthInBytes+1)),
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
				},
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
		// Test failure when there is too much custom attributes
		{
			input: &Input{
				Payload:          jwtPayloadBody,
				AttributeMapping: generateAttrMap(MaximumCustomAttributes + 1),
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
		// Test failure when timestamp is instantiated with an integer
		{
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleSubject:                           attribute.GetAssertionName("sub"),
					attribute.GetAttributeName("timestamp"): "string(timestamp(int(assertion.iat)))",
				},
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
	}

	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			c := Compiler{
				Input:    tc.input,
				Provider: tc.provider,
			}
			_, err := c.Run()

			if tc.expected.IsError {
				if err == nil {
					t.Fatalf("Run(%v) = %s, expected an error", c.Input, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Run(%v) = %s, expected no error", c.Input, err)
				}
			}
		})
	}
}

func TestAttributeCondition(t *testing.T) {
	type Expected struct {
		IsError   bool
		ErrorType error
	}
	tests := []struct {
		input    *Input
		provider provider.Provider
		expected *Expected
	}{
		{
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
					GoogleGroups:  attribute.GetAssertionName("groups"),
				},
				AttributeCondition: `google.subject == "1234567890" && "group1" in google.groups && assertion.is_admin == true`,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: false,
			},
		},
		{
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
				},
				AttributeCondition: `assertion.is_admin == true && "group1" in assertion.groups`,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: false,
			},
		},
		{
			input: &Input{
				Payload:            `{"sub": "1234567890", "is_admin": true, "iat": 1516239022}`,
				AttributeMapping:   map[string]string{GoogleSubject: attribute.GetAssertionName("sub")},
				AttributeCondition: `assertion.is_admin == false`,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError:   true,
				ErrorType: ErrAttrConditionFailed,
			},
		},
		{
			input: &Input{
				Payload:            `{"sub": "1234567890", "is_admin": true, "iat": 1516239022}`,
				AttributeMapping:   map[string]string{GoogleSubject: attribute.GetAssertionName("sub")},
				AttributeCondition: `assertion.is_admin == false`,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError:   true,
				ErrorType: ErrAttrConditionFailed,
			},
		},
		{
			input: &Input{
				Payload:            `{"sub": "1234567890", "is_admin": true, "iat": 1516239022}`,
				AttributeMapping:   map[string]string{GoogleSubject: attribute.GetAssertionName("sub")},
				AttributeCondition: `'string'`,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError:   true,
				ErrorType: ErrAttrConditionFailed,
			},
		},
		{
			input: &Input{
				Payload:            `{"sub": "1234567890", "is_admin": true, "iat": 1516239022}`,
				AttributeMapping:   map[string]string{GoogleSubject: attribute.GetAssertionName("sub")},
				AttributeCondition: invalidCelExpr,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
		// Test failure when timestamp is instantiated with an integer
		{
			input: &Input{
				Payload:            jwtPayloadBody,
				AttributeMapping:   map[string]string{GoogleSubject: attribute.GetAssertionName("sub")},
				AttributeCondition: `timestamp(int(assertion.iat)) == timestamp(int(assertion.iat))`,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				IsError: true,
			},
		},
	}

	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			c := Compiler{
				Input:    tc.input,
				Provider: tc.provider,
			}
			_, err := c.Run()

			if tc.expected.IsError {
				if err == nil {
					t.Fatalf("Run(%v) = %s, expected an error", c.Input, err)
				} else if tc.expected.ErrorType != nil && err != tc.expected.ErrorType { //nolint:errorlint
					t.Fatalf("Run(%v) = %s, expected %s", c.Input, err, tc.expected.ErrorType)
				}
			} else {
				if err != nil {
					t.Fatalf("Run(%v) = %s, expected no error", c.Input, err)
				}
			}
		})
	}
}
