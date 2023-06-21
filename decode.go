package applereceipt

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"

	"go.mozilla.org/pkcs7"
)

//go:generate go run gen.go

type ReceiptAttribute struct {
	Type    int
	Version int
	Value   []byte
}

// DecodeBase64 decodes given base64-encoded receipt into AppReceipt type.
// It verifies the payload with given certPool. certPool can be nil to skip verification, but be careful that users can pass a forged receipt.
func DecodeBase64(b64Receipt string, certPool *x509.CertPool) (AppReceipt, error) {
	receipt := make([]byte, base64.StdEncoding.DecodedLen(len(b64Receipt)))
	_, err := base64.StdEncoding.Decode(receipt, []byte(b64Receipt))
	if err != nil {
		return AppReceipt{}, err
	}
	return Decode(receipt, certPool)
}

// Decode decodes given receipt binary into AppReceipt type.
// It verifies the payload with given certPool. certPool can be nil to skip verification, but be careful that users can pass a forged receipt.
func Decode(receipt []byte, certPool *x509.CertPool) (AppReceipt, error) {
	p, err := pkcs7.Parse(receipt)
	if err != nil {
		return AppReceipt{}, err
	}

	var attrs []ReceiptAttribute
	if _, err := asn1.UnmarshalWithParams(p.Content, &attrs, "set"); err != nil {
		return AppReceipt{}, err
	}
	appReceipt, err := newAppReceipt(attrs)
	if err != nil {
		return AppReceipt{}, err
	}

	if err := p.VerifyWithChainAtTime(certPool, appReceipt.ReceiptCreationDate); err != nil {
		return AppReceipt{}, err
	}

	return appReceipt, nil
}
