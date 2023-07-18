# go-applereceipt [![Test](https://github.com/devsisters/go-applereceipt/actions/workflows/test.yaml/badge.svg)](https://github.com/devsisters/go-applereceipt/actions/workflows/test.yaml) [![Go Reference](https://pkg.go.dev/badge/github.com/devsisters/go-applereceipt.svg)](https://pkg.go.dev/github.com/devsisters/go-applereceipt)

Go library to parse and verify Apple App Store receipts, locally on the server.

> The verifyReceipt endpoint is deprecated. To validate receipts on your server, follow the steps in [Validating receipts on the device](https://developer.apple.com/documentation/appstorereceipts/validating_receipts_on_the_device) on your server.
> — https://developer.apple.com/documentation/appstorereceipts/verifyreceipt

This library implements the PKCS#7 signature verification and ASN.1 parsing of the payload locally on the server, without the need to call Apple's `verifyReceipt` endpoint. The parsed receipt is filled in a concrete `AppReceipt` struct, which is auto-generated based on [the documented fields](https://developer.apple.com/library/archive/releasenotes/General/ValidateAppStoreReceipt/Chapters/ReceiptFields.html) published by Apple.

Note that at this moment, receipts from Apple is signed using SHA1-RSA, which has [removed support since Go 1.18](https://go.dev/issue/41682) and requires `GODEBUG=x509sha1=1` to verify receipts. SHA-256 signatures will be available from receipts since August 14, 2023. ([TN3138](https://developer.apple.com/documentation/technotes/tn3138-handling-app-store-receipt-signing-certificate-changes))

## Installation

```sh
go get github.com/devsisters/go-applereceipt
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/devsisters/go-applereceipt"
	"github.com/devsisters/go-applereceipt/applepki"
)

func main() {
	receipt, err := applereceipt.DecodeBase64("MIIT...", applepki.CertPool())
	if err != nil {
		panic(err)
	}

	fmt.Println(receipt.BundleIdentifier)
}
```
