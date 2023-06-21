package applepki

import (
	"crypto/x509"
	"embed"
	"sync"
)

//go:generate go run gen.go

//go:embed certs/*.cer
var certs embed.FS
var pool *x509.CertPool
var poolOnce sync.Once

func CertPool() *x509.CertPool {
	poolOnce.Do(func() {
		pool = x509.NewCertPool()
		entries, err := certs.ReadDir("certs")
		if err != nil {
			return
		}
		for _, entry := range entries {
			if !entry.IsDir() && entry.Type().IsRegular() {
				cert, err := certs.ReadFile("certs/" + entry.Name())
				if err != nil {
					continue
				}
				isPem := pool.AppendCertsFromPEM(cert)
				if !isPem {
					if cer, err := x509.ParseCertificate(cert); err == nil {
						pool.AddCert(cer)
					}
				}
			}
		}
	})
	return pool
}
