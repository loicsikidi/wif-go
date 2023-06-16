package compiler

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/loicsikidi/wif-go/pkg/compiler/provider"
	"github.com/loicsikidi/wif-go/pkg/compiler/provider/attribute"
	"github.com/loicsikidi/wif-go/pkg/compiler/provider/oidc"
	"github.com/loicsikidi/wif-go/pkg/compiler/provider/saml"
	samlType "github.com/loicsikidi/wif-go/pkg/saml"
)

type Payload struct {
	Sub     string   `json:"sub"`
	IsAdmin bool     `json:"is_admin"`
	Groups  []string `json:"groups"`
}

func (p *Payload) GetMap() map[string]any {
	return map[string]any{"sub": p.Sub, "is_admin": p.IsAdmin, "groups": p.Groups}
}

const (
	invalidCelExpr = "invalid.cel.expr"
	jwtPayloadBody = `{"sub": "1234567890", "is_admin": true, "iat": 1683438895, "groups": ["group1", "group2"]}`
)

var (
	defaultPayload     = Payload{Sub: "1234567890", IsAdmin: true, Groups: []string{"group1", "group2"}}
	defaultJwtPayload  = generateJWT(defaultPayload.GetMap())
	defaultSamlPayload = generationSAML()
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

func generateJWT(input map[string]any) string {
	// Set the claims
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	for k, v := range input {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte("secret")
	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func generationSAML() string {
	response := samlType.Response{
		Assertion: &samlType.Assertion{
			ID:           "_110140b8f596cfc7772e93b1e2d69f9",
			IssueInstant: time.Now(),
			Version:      "2.0",
			Subject: &samlType.Subject{
				NameID: &samlType.NameID{
					Value: "1234567890",
				},
			},
			AttributeStatement: &samlType.AttributeStatement{
				Attributes: []samlType.Attribute{
					{
						Name:       "groups",
						NameFormat: "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
						Values: []samlType.AttributeValue{
							{
								Value: "group1",
								Type:  "xs:string",
							},
							{
								Value: "group2",
								Type:  "xs:string",
							},
						},
					},
				},
			},
		},
	}
	xmlBytes, _ := xml.Marshal(response)
	return string(xmlBytes)
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
	tests := []struct {
		name     string
		input    *Input
		provider string
		wantErr  bool
	}{
		{
			name: "test success saml provider with xml payload",
			input: &Input{
				Payload: defaultSamlPayload,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("subject"),
					GoogleGroups:  attribute.GetAssertionName("attributes['groups']"),
				},
			},
			provider: "saml",
			wantErr:  false,
		},
		{
			name: "test success saml provider with ambient feature",
			input: &Input{
				Payload: defaultSamlPayload,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("subject"),
					GoogleGroups:  attribute.GetAssertionName("attributes['groups']"),
				},
			},
			provider: "ambient",
			wantErr:  false,
		},
		{
			name: "test success oidc provider with json payload",
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
					GoogleGroups:  attribute.GetAssertionName("groups"),
				},
			},
			provider: "oidc",
			wantErr:  false,
		},
		{
			name: "test success oidc provider with jwt payload",
			input: &Input{
				Payload: defaultJwtPayload,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
					GoogleGroups:  attribute.GetAssertionName("groups"),
				},
			},
			provider: "oidc",
			wantErr:  false,
		},
		{
			name: "test success oidc provider with ambient feature",
			input: &Input{
				Payload: defaultJwtPayload,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
					GoogleGroups:  attribute.GetAssertionName("groups"),
				},
			},
			provider: "ambient",
			wantErr:  false,
		},
		{
			name: "test failure if the mapped value is not a string",
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleSubject: "true",
				},
			},
			provider: "oidc",
			wantErr:  true,
		},
		{
			name: "test failure if google.groups is not a list",
			input: &Input{
				Payload: `{"sub":"1234567890", "groups": "group1"}`,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
					GoogleGroups:  attribute.GetAssertionName("groups"),
				},
			},
			provider: "oidc",
			wantErr:  true,
		},
		{
			name: "test failure if any element of google.groups is not a string",
			input: &Input{
				Payload: `{"sub":"1234567890", "groups": ["group1", "group2", 1]}`,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
					GoogleGroups:  attribute.GetAssertionName("groups"),
				},
			},
			provider: "oidc",
			wantErr:  true,
		},
		{
			name: "test failure if the payload is invalid JSON",
			input: &Input{
				Payload: `{sub: "1234567890"}`,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
				},
			},
			provider: "oidc",
			wantErr:  true,
		},
		{
			name: "test failure if the mapped key is unauthorized",
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					"invalid_key": attribute.GetAssertionName("sub"),
				},
			},
			provider: "oidc",
			wantErr:  true,
		},
		{
			name: "test failure if missing subject attribute",
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleGroups: attribute.GetAssertionName("groups"),
				},
			},
			provider: "oidc",
			wantErr:  true,
		},
		{
			name: "test failure if the attribute subject value is too big",
			input: &Input{
				Payload: fmt.Sprintf(`{"sub": "%s"}`, generateStr(MaximumSubjectLengthInBytes+1)),
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
				},
			},
			provider: "oidc",
			wantErr:  true,
		},
		{
			name: "test failure if there is too much custom attributes",
			input: &Input{
				Payload:          jwtPayloadBody,
				AttributeMapping: generateAttrMap(MaximumCustomAttributes + 1),
			},
			provider: "oidc",
			wantErr:  true,
		},
		{
			name: "test failure if timestamp is instantiated with an integer",
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleSubject:                           attribute.GetAssertionName("sub"),
					attribute.GetAttributeName("timestamp"): "string(timestamp(int(assertion.iat)))",
				},
			},
			provider: "oidc",
			wantErr:  true,
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			var provider provider.Provider
			switch tc.provider {
			case "oidc":
				provider = &oidc.Provider{}
			case "saml":
				provider = &saml.Provider{}
			case "ambient":
			default:
				t.Fatalf("unknown provider %s", tc.provider)
			}
			c := Compiler{
				Input:    tc.input,
				Provider: provider,
			}
			_, err := c.Run()

			if tc.wantErr {
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
		wantErr   bool
		ErrorType error
	}
	tests := []struct {
		name     string
		input    *Input
		provider provider.Provider
		expected *Expected
	}{
		{
			name: "test success if condition return a true boolean",
			input: &Input{
				Payload: jwtPayloadBody,
				AttributeMapping: map[string]string{
					GoogleSubject: attribute.GetAssertionName("sub"),
					GoogleGroups:  attribute.GetAssertionName("groups"),
				},
				AttributeCondition: `google.subject == "1234567890" && "group1" in google.groups && assertion.is_admin == true && "group1" in assertion.groups`,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				wantErr: false,
			},
		},
		{
			name: "test failure if condition return a false boolean",
			input: &Input{
				Payload:            `{"sub": "1234567890", "is_admin": true, "iat": 1516239022}`,
				AttributeMapping:   map[string]string{GoogleSubject: attribute.GetAssertionName("sub")},
				AttributeCondition: `assertion.is_admin == false`,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				wantErr:   true,
				ErrorType: ErrAttrConditionFailed,
			},
		},
		{
			name: "test failure if condition is not a boolean",
			input: &Input{
				Payload:            `{"sub": "1234567890", "is_admin": true, "iat": 1516239022}`,
				AttributeMapping:   map[string]string{GoogleSubject: attribute.GetAssertionName("sub")},
				AttributeCondition: `'string'`,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				wantErr:   true,
				ErrorType: ErrAttrConditionFailed,
			},
		},
		{
			name: "test failure if there is invalid cel expression",
			input: &Input{
				Payload:            `{"sub": "1234567890", "is_admin": true, "iat": 1516239022}`,
				AttributeMapping:   map[string]string{GoogleSubject: attribute.GetAssertionName("sub")},
				AttributeCondition: invalidCelExpr,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				wantErr: true,
			},
		},
		{
			name: "test failure if timestamp is instantiated with an integer",
			input: &Input{
				Payload:            jwtPayloadBody,
				AttributeMapping:   map[string]string{GoogleSubject: attribute.GetAssertionName("sub")},
				AttributeCondition: `timestamp(int(assertion.iat)) == timestamp(int(assertion.iat))`,
			},
			provider: &oidc.Provider{},
			expected: &Expected{
				wantErr: true,
			},
		},
	}

	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			c := Compiler{
				Input:    tc.input,
				Provider: tc.provider,
			}
			_, err := c.Run()

			if tc.expected.wantErr {
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
