package saml

import (
	"encoding/xml"
	"time"

	"github.com/loicsikidi/wif-go/pkg/saml/xmlsec"
)

const (
	SubjectStr    string = "subject"
	AttributesStr string = "attributes"
)

// Issuer represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type Issuer struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`
	Format  string   `xml:",attr"`
	Value   string   `xml:",chardata"`
}

// NameIDPolicy represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
// Also refer to Azure docs for their IdP supported values: https://msdn.microsoft.com/en-us/library/azure/dn195589.aspx
type NameIDPolicy struct {
	XMLName xml.Name

	// Optional attributes
	//

	// A Boolean value used to indicate whether the identity provider is allowed, in the course of fulfilling the
	// request, to create a new identifier to represent the principal. Defaults to "false". When "false", the
	// requester constrains the identity provider to only issue an assertion to it if an acceptable identifier for
	// the principal has already been established. Note that this does not prevent the identity provider from
	// creating such identifiers outside the context of this specific request (for example, in advance for a
	// large number of principals)
	AllowCreate bool `xml:",attr"`

	// Specifies the URI reference corresponding to a name identifier format defined in this or another
	// specification (see Section 8.3 for examples). The additional value of
	// urn:oasis:names:tc:SAML:2.0:nameid-format:encrypted is defined specifically for use
	// within this attribute to indicate a request that the resulting identifier be encrypted
	Format string `xml:",attr"`
}

// Response represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf section 3.2
type Response struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol Response"`

	// Required attributes
	//

	// An identifier for the request.
	// The values of the ID attribute in a request and the InResponseTo
	// attribute in the corresponding response MUST match.
	ID string `xml:",attr"`

	// The version of this request.
	// Only version 2.0 is supported by pressly/saml
	Version string `xml:",attr"`

	// The time instant of issue of the request. The time value is encoded in UTC
	IssueInstant time.Time `xml:",attr"`

	// A code representing the status of the corresponding reques
	Status *Status

	// Optional attributes
	//

	// A URI reference indicating the address to which this request has been sent. This is useful to prevent
	// malicious forwarding of requests to unintended recipients, a protection that is required by some
	// protocol bindings. If it is present, the actual recipient MUST check that the URI reference identifies the
	// location at which the message was received. If it does not, the request MUST be discarded. Some
	// protocol bindings may require the use of this attribute
	Destination string `xml:",attr"`

	// An XML Signature that authenticates the requester and provides message integrity
	Signature *xmlsec.Signature

	// A reference to the identifier of the request to which the response corresponds, if any. If the response
	// is not generated in response to a request, or if the ID attribute value of a request cannot be
	// determined (for example, the request is malformed), then this attribute MUST NOT be present.
	// Otherwise, it MUST be present and its value MUST match the value of the corresponding request's
	// ID attribute.
	InResponseTo string `xml:",attr"`

	// Identifies the entity that generated the request message
	// By default, the value of the <Issuer> element is a URI of no more than 1024 characters.
	// Changes from SAML version 1 to 2
	// An <Issuer> element can now be present on requests and responses (in addition to appearing on assertions).
	Issuer *Issuer

	EncryptedAssertion *EncryptedAssertion

	Assertion *Assertion
}

// GetAttributes returns the attributes of the SAML response
func (r *Response) GetMapAttributes() map[string]any {
	m := map[string]any{}
	m[SubjectStr] = r.Assertion.Subject.NameID.Value
	m[AttributesStr] = make(map[string]any)
	for _, attr := range r.Assertion.AttributeStatement.Attributes {
		attrMap := map[string]any{}
		attrMap[attr.Name] = make([]any, 0)
		for _, value := range attr.Values {
			attrMap[attr.Name] = append(attrMap[attr.Name].([]any), value.Value)
		}
		m[AttributesStr] = attrMap
	}
	return m
}

// Status represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type Status struct {
	XMLName    xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol Status"`
	StatusCode StatusCode
}

// StatusCode represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type StatusCode struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:protocol StatusCode"`
	Value   string   `xml:",attr"`
}

// StatusSuccess is the value of a StatusCode element when the authentication succeeds.
// (nominally a constant, except for testing)
var StatusSuccess = "urn:oasis:names:tc:SAML:2.0:status:Success"

// EncryptedAssertion represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type EncryptedAssertion struct {
	XMLName       xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion EncryptedAssertion"`
	Assertion     *Assertion
	EncryptedData []byte `xml:",innerxml"`
}

// Assertion represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type Assertion struct {
	XMLName            xml.Name  `xml:"urn:oasis:names:tc:SAML:2.0:assertion Assertion"`
	ID                 string    `xml:",attr"`
	IssueInstant       time.Time `xml:",attr"`
	Version            string    `xml:",attr"`
	Issuer             *Issuer
	Signature          *xmlsec.Signature
	Subject            *Subject
	Conditions         *Conditions
	AuthnStatement     *AuthnStatement
	AttributeStatement *AttributeStatement
}

// Subject represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type Subject struct {
	XMLName             xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion Subject"`
	NameID              *NameID
	SubjectConfirmation *SubjectConfirmation
}

// NameID represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type NameID struct {
	Format          string `xml:",attr"`
	NameQualifier   string `xml:",attr"`
	SPNameQualifier string `xml:",attr"`
	Value           string `xml:",chardata"`
}

// SubjectConfirmation represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type SubjectConfirmation struct {
	Method                  string `xml:",attr"`
	SubjectConfirmationData SubjectConfirmationData
}

// SubjectConfirmationData represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type SubjectConfirmationData struct {
	Address      string    `xml:",attr"`
	InResponseTo string    `xml:",attr"`
	NotOnOrAfter time.Time `xml:",attr"`
	Recipient    string    `xml:",attr"`
}

// Conditions represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type Conditions struct {
	NotBefore           time.Time `xml:",attr"`
	NotOnOrAfter        time.Time `xml:",attr"`
	AudienceRestriction *AudienceRestriction
}

// AudienceRestriction represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type AudienceRestriction struct {
	Audience *Audience
}

// Audience represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type Audience struct {
	Value string `xml:",chardata"`
}

// AuthnStatement represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type AuthnStatement struct {
	AuthnInstant    time.Time `xml:",attr"`
	SessionIndex    string    `xml:",attr"`
	SubjectLocality SubjectLocality
	AuthnContext    AuthnContext
}

// SubjectLocality represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type SubjectLocality struct {
	Address string `xml:",attr"`
}

// AuthnContext represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type AuthnContext struct {
	AuthnContextClassRef *AuthnContextClassRef
}

// AuthnContextClassRef represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type AuthnContextClassRef struct {
	Value string `xml:",chardata"`
}

// AttributeStatement represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type AttributeStatement struct {
	Attributes []Attribute `xml:"Attribute"`
}

// Attribute represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type Attribute struct {
	FriendlyName string           `xml:",attr"`
	Name         string           `xml:",attr"`
	NameFormat   string           `xml:",attr"`
	Values       []AttributeValue `xml:"AttributeValue"`
}

// AttributeValue represents the SAML object of the same name.
//
// See http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
type AttributeValue struct {
	Type   string `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
	Value  string `xml:",chardata"`
	NameID *NameID
}
