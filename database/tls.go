package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// buildTLSConfig creates a tls.Config based on the requested mode and certificate paths.
func buildTLSConfig(tlsMode, serverName, rootCertPath, clientCertPath, clientKeyPath string) (*tls.Config, error) {
	if tlsMode == "disable" || tlsMode == "" {
		return nil, nil
	}

	tlsConfig := &tls.Config{}

	// Set server name for hostname verification in "verify-full" mode.
	if tlsMode == "verify-full" {
		tlsConfig.ServerName = serverName
	}

	// In "require" mode, we don't verify the server's cert, so we skip the rest.
	if tlsMode == "require" {
		tlsConfig.InsecureSkipVerify = true
		return tlsConfig, nil
	}

	// Load custom root CA if provided, otherwise use system's trust store.
	if rootCertPath != "" {
		caCert, err := os.ReadFile(rootCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read root certificate: %w", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	// Load client certificate and key for mTLS if both are provided.
	if clientCertPath != "" && clientKeyPath != "" {
		clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load client key pair: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	} else if clientCertPath != "" || clientKeyPath != "" {
		return nil, fmt.Errorf("both client_cert_path and client_key_path must be provided for mTLS")
	}

	return tlsConfig, nil
}
