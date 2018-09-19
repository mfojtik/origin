package assets

import (
	"crypto/rsa"
	"crypto/x509"
	"net"

	"github.com/google/uuid"
	"github.com/openshift/origin/pkg/oc/clusterup/coreinstall/assets/tlsutil"
)

func (r *RenderOptions) newKeyAndCert(caCert *x509.Certificate, caPrivKey *rsa.PrivateKey, commonName string, addrs []string) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tlsutil.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	var altNames tlsutil.AltNames
	for _, addr := range addrs {
		if ip := net.ParseIP(addr); ip != nil {
			altNames.IPs = append(altNames.IPs, ip)
		} else {
			altNames.DNSNames = append(altNames.DNSNames, addr)
		}
	}
	config := tlsutil.CertConfig{
		CommonName:   commonName,
		Organization: []string{"etcd"},
		AltNames:     altNames,
	}
	cert, err := tlsutil.NewSignedCertificate(config, key, caCert, caPrivKey)
	if err != nil {
		return nil, nil, err
	}
	return key, cert, err
}

func (r *RenderOptions) newCACert() (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tlsutil.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	config := tlsutil.CertConfig{
		CommonName:         r.CommonName,
		Organization:       []string{uuid.New()},
		OrganizationalUnit: r.OrganizationalUnitNames,
	}

	cert, err := tlsutil.NewSelfSignedCACertificate(config, key)
	if err != nil {
		return nil, nil, err
	}

	return key, cert, err
}
