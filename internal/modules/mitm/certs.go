package mitm

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	caCommonName     = "OMEG4MON Local MITM Root CA"
	caValidityDays   = 3650
	leafValidityDays = 7
	maxLeafCacheSize = 256
)

// ErrCANotFound indicates the Root CA certificate has not been generated yet.
var ErrCANotFound = errors.New("mitm: root CA certificate not generated")

// CAInfo describes the public metadata of the local MITM Root CA.
type CAInfo struct {
	Exists      bool   `json:"exists"`
	CertPath    string `json:"cert_path"`
	Fingerprint string `json:"fingerprint"`
	Subject     string `json:"subject"`
	NotBefore   string `json:"not_before"`
	NotAfter    string `json:"not_after"`
}

type cachedLeaf struct {
	cert      *tls.Certificate
	expiresAt time.Time
}

type certManager struct {
	dir string

	mu     sync.RWMutex
	caCert *x509.Certificate
	caKey  *ecdsa.PrivateKey
	caPEM  []byte
	leaves map[string]cachedLeaf
}

func newCertManager(dataDir string) *certManager {
	dir := filepath.Join(dataDir, "mitm")
	return &certManager{
		dir:    dir,
		leaves: make(map[string]cachedLeaf),
	}
}

func (c *certManager) certPath() string { return filepath.Join(c.dir, "rootCA.crt") }
func (c *certManager) keyPath() string  { return filepath.Join(c.dir, "rootCA.key") }

// loadFromDisk loads an existing Root CA if present on disk.
func (c *certManager) loadFromDisk() error {
	certPEM, err := os.ReadFile(c.certPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read root CA cert: %w", err)
	}
	keyPEM, err := os.ReadFile(c.keyPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read root CA key: %w", err)
	}

	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return errors.New("mitm: invalid root CA certificate PEM")
	}
	parsedCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return fmt.Errorf("parse root CA cert: %w", err)
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return errors.New("mitm: invalid root CA key PEM")
	}
	parsedKey, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("parse root CA key: %w", err)
	}

	c.mu.Lock()
	c.caCert = parsedCert
	c.caKey = parsedKey
	c.caPEM = certPEM
	c.leaves = make(map[string]cachedLeaf)
	c.mu.Unlock()
	return nil
}

// ensureCA returns existing Root CA metadata or generates a new ECDSA P-256 Root CA.
func (c *certManager) ensureCA(force bool) (CAInfo, error) {
	c.mu.RLock()
	hasCA := c.caCert != nil && c.caKey != nil
	c.mu.RUnlock()
	if hasCA && !force {
		return c.info(), nil
	}

	if err := os.MkdirAll(c.dir, 0o700); err != nil {
		return CAInfo{}, fmt.Errorf("create mitm cert dir: %w", err)
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return CAInfo{}, fmt.Errorf("generate root CA key: %w", err)
	}

	serial, err := randomSerial()
	if err != nil {
		return CAInfo{}, err
	}
	now := time.Now().UTC().Add(-5 * time.Minute)
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   caCommonName,
			Organization: []string{"OMEG4MON ROUTER"},
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(0, 0, caValidityDays),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &priv.PublicKey, priv)
	if err != nil {
		return CAInfo{}, fmt.Errorf("create root CA cert: %w", err)
	}
	parsedCert, err := x509.ParseCertificate(der)
	if err != nil {
		return CAInfo{}, fmt.Errorf("parse generated root CA cert: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return CAInfo{}, fmt.Errorf("marshal root CA key: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	if err := writeAtomicFile(c.certPath(), certPEM, 0o644); err != nil {
		return CAInfo{}, fmt.Errorf("write root CA cert: %w", err)
	}
	if err := writeAtomicFile(c.keyPath(), keyPEM, 0o600); err != nil {
		return CAInfo{}, fmt.Errorf("write root CA key: %w", err)
	}

	c.mu.Lock()
	c.caCert = parsedCert
	c.caKey = priv
	c.caPEM = certPEM
	c.leaves = make(map[string]cachedLeaf)
	c.mu.Unlock()

	return c.info(), nil
}

func (c *certManager) info() CAInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := CAInfo{
		Exists:   c.caCert != nil,
		CertPath: c.certPath(),
	}
	if c.caCert == nil {
		return out
	}
	sum := sha256.Sum256(c.caCert.Raw)
	out.Fingerprint = hex.EncodeToString(sum[:])
	out.Subject = c.caCert.Subject.String()
	out.NotBefore = c.caCert.NotBefore.UTC().Format(time.RFC3339)
	out.NotAfter = c.caCert.NotAfter.UTC().Format(time.RFC3339)
	return out
}

func (c *certManager) downloadPEM() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.caPEM) == 0 {
		return nil, ErrCANotFound
	}
	cp := make([]byte, len(c.caPEM))
	copy(cp, c.caPEM)
	return cp, nil
}

func (c *certManager) getCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	host := "localhost"
	if hello != nil && strings.TrimSpace(hello.ServerName) != "" {
		host = strings.ToLower(strings.TrimSpace(hello.ServerName))
	}
	if h, _, err := net.SplitHostPort(host); err == nil && h != "" {
		host = h
	}

	now := time.Now().UTC()
	c.mu.RLock()
	if cached, ok := c.leaves[host]; ok && cached.expiresAt.After(now.Add(time.Hour)) {
		c.mu.RUnlock()
		return cached.cert, nil
	}
	caCert := c.caCert
	caKey := c.caKey
	c.mu.RUnlock()

	if caCert == nil || caKey == nil {
		return nil, ErrCANotFound
	}

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate leaf key: %w", err)
	}
	serial, err := randomSerial()
	if err != nil {
		return nil, err
	}

	expiresAt := now.AddDate(0, 0, leafValidityDays)
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: host,
		},
		NotBefore:   now.Add(-5 * time.Minute),
		NotAfter:    expiresAt,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	if ip := net.ParseIP(host); ip != nil {
		tmpl.IPAddresses = []net.IP{ip}
	} else {
		tmpl.DNSNames = []string{host}
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, &leafKey.PublicKey, caKey)
	if err != nil {
		return nil, fmt.Errorf("sign leaf cert for %q: %w", host, err)
	}
	parsedLeaf, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, fmt.Errorf("parse leaf cert: %w", err)
	}
	tlsCert := &tls.Certificate{
		Certificate: [][]byte{der, caCert.Raw},
		PrivateKey:  leafKey,
		Leaf:        parsedLeaf,
	}

	c.mu.Lock()
	if len(c.leaves) >= maxLeafCacheSize {
		clear(c.leaves)
	}
	c.leaves[host] = cachedLeaf{cert: tlsCert, expiresAt: expiresAt}
	c.mu.Unlock()

	return tlsCert, nil
}

func randomSerial() (*big.Int, error) {
	limit := new(big.Int).Lsh(big.NewInt(1), 128)
	n, err := rand.Int(rand.Reader, limit)
	if err != nil {
		return nil, fmt.Errorf("generate cert serial: %w", err)
	}
	return n, nil
}

func writeAtomicFile(path string, data []byte, perm os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, perm); err != nil {
		return err
	}
	if err := os.Chmod(tmp, perm); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, path)
}
