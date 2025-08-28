package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

// buildTLSConfig creates a tls.Config based on the requested mode.
// It returns nil for "disable" mode.
func buildTLSConfig(tlsMode string, serverName string) (*tls.Config, error) {
	switch tlsMode {
	case "disable", "":
		return nil, nil
	case "require":
		// This is insecure, but matches the user's request for "require" mode.
		return &tls.Config{InsecureSkipVerify: true}, nil
	case "verify-ca":
		// Get the system's CA pool.
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			// In some environments (like Windows) this can be nil.
			// Fallback to an empty pool.
			rootCAs = x509.NewCertPool()
		}
		return &tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            rootCAs,
		}, nil
	case "verify-full":
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}
		return &tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            rootCAs,
			ServerName:         serverName,
		}, nil
	default:
		return nil, fmt.Errorf("invalid tls_mode: %s", tlsMode)
	}
}
