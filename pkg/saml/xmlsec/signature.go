package xmlsec

import (
	"encoding/xml"
)

// Method is part of Signature.
type Method struct {
	Algorithm string `xml:",attr"`
}

// Signature is a model for the Signature object specified by XMLDSIG.
type Signature struct {
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Signature"`

	CanonicalizationMethod Method             `xml:"SignedInfo>CanonicalizationMethod"`
	SignatureMethod        Method             `xml:"SignedInfo>SignatureMethod"`
	Reference              Reference          `xml:"SignedInfo>Reference"`
	SignatureValue         string             `xml:"SignatureValue"`
	KeyName                string             `xml:"KeyInfo>KeyName,omitempty"`
	X509Certificate        *SignatureX509Data `xml:"KeyInfo>X509Data,omitempty"`
}

type Reference struct {
	URI          string   `xml:"URI,attr,omitempty"`
	Transforms   []Method `xml:"Transforms>Transform"`
	DigestMethod Method   `xml:"DigestMethod"`
	DigestValue  string   `xml:"DigestValue"`
}

// SignatureX509Data represents the <X509Data> element of <Signature>
type SignatureX509Data struct {
	X509Certificate string `xml:"X509Certificate,omitempty"`
}
