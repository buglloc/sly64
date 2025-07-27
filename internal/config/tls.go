package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"github.com/buglloc/certifi"

	"github.com/buglloc/sly64/v2/internal/config/configpb"
)

func NewTLSConfig(cfg *configpb.TLS) (*tls.Config, error) {
	tlsCfg := tls.Config{
		RootCAs: certifi.NewCertPool(),
	}
	if cfg == nil {
		return &tlsCfg, nil
	}

	tlsCfg.ServerName = cfg.ServerName
	tlsCfg.InsecureSkipVerify = cfg.InsecureSkipVerify
	if len(cfg.CaCert) > 0 {
		cas, err := newCertPool(cfg.CaCert)
		if err != nil {
			return nil, fmt.Errorf("create root CAs: %w", err)
		}

		tlsCfg.RootCAs = cas
	}

	return &tlsCfg, nil
}

func newCertPool(caCert string) (*x509.CertPool, error) {
	in, err := os.ReadFile(caCert)
	if err != nil {
		return nil, fmt.Errorf("read cacert: %w", err)
	}

	certs, err := certifi.ParseCertificates(in)
	if err != nil {
		return nil, fmt.Errorf("parse cacert: %w", err)
	}

	out := x509.NewCertPool()
	for _, cert := range certs {
		out.AddCert(cert)
	}

	return out, nil
}

func patchTLSConfig(cfg *configpb.TLS, path string) error {
	if cfg == nil {
		return nil
	}

	if len(cfg.CaCert) == 0 {
		return nil
	}

	if filepath.IsAbs(cfg.CaCert) {
		return nil
	}

	cfg.CaCert = filepath.Join(filepath.Dir(path), cfg.CaCert)
	return nil
}
