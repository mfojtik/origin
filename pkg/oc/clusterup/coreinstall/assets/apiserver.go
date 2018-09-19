package assets

import (
	"crypto/rsa"
	"crypto/x509"

	assetslib "github.com/openshift/library-go/pkg/assets"
	"github.com/openshift/origin/pkg/oc/clusterup/coreinstall/assets/tlsutil"
)

const (
	AssetPathCAKey                 = ""
	AssetPathCACert                = ""
	AssetPathAPIServerKey          = ""
	AssetPathAPIServerCert         = ""
	AssetPathServiceAccountPrivKey = ""
	AssetPathServiceAccountPubKey  = ""
)

func (r *RenderOptions) newAPIKeyAndCert(caCert *x509.Certificate, caPrivKey *rsa.PrivateKey, altNames tlsutil.AltNames) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tlsutil.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	altNames.DNSNames = append(altNames.DNSNames, []string{
		"kubernetes",
		"kubernetes.default",
		"kubernetes.default.svc",
		"kubernetes.default.svc.cluster.local",
	}...)

	config := tlsutil.CertConfig{
		CommonName:   "kube-apiserver",
		Organization: []string{"kube-master"},
		AltNames:     altNames,
	}
	cert, err := tlsutil.NewSignedCertificate(config, key, caCert, caPrivKey)
	if err != nil {
		return nil, nil, err
	}
	return key, cert, err
}

func (r *RenderOptions) newTLSAssets(caCert *x509.Certificate, caPrivKey *rsa.PrivateKey, altNames tlsutil.AltNames) ([]assetslib.Asset, error) {
	var (
		assets []assetslib.Asset
		err    error
	)

	apiKey, apiCert, err := r.newAPIKeyAndCert(caCert, caPrivKey, altNames)
	if err != nil {
		return assets, err
	}

	saPrivKey, err := tlsutil.NewPrivateKey()
	if err != nil {
		return assets, err
	}

	saPubKey, err := tlsutil.EncodePublicKeyPEM(&saPrivKey.PublicKey)
	if err != nil {
		return assets, err
	}

	assets = append(assets, []assetslib.Asset{
		{Name: AssetPathCAKey, Data: tlsutil.EncodePrivateKeyPEM(caPrivKey)},
		{Name: AssetPathCACert, Data: tlsutil.EncodeCertificatePEM(caCert)},
		{Name: AssetPathAPIServerKey, Data: tlsutil.EncodePrivateKeyPEM(apiKey)},
		{Name: AssetPathAPIServerCert, Data: tlsutil.EncodeCertificatePEM(apiCert)},
		{Name: AssetPathServiceAccountPrivKey, Data: tlsutil.EncodePrivateKeyPEM(saPrivKey)},
		{Name: AssetPathServiceAccountPubKey, Data: saPubKey},
	}...)
	return assets, nil
}
