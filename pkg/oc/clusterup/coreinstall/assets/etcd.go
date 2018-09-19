package assets

import (
	"crypto/rsa"
	"crypto/x509"
	"net/url"

	assetslib "github.com/openshift/library-go/pkg/assets"

	"github.com/openshift/origin/pkg/oc/clusterup/coreinstall/assets/tlsutil"
)

const (
	AssetPathEtcdClientCA   = ""
	AssetPathEtcdClientKey  = ""
	AssetPathEtcdClientCert = ""

	AssetPathEtcdPeerCA     = ""
	AssetPathEtcdPeerCert   = ""
	AssetPathEtcdServerCA   = ""
	AssetPathEtcdServerKey  = ""
	AssetPathEtcdServerCert = ""
)

func (r *RenderOptions) newEtcdTLSAssets(caCert *x509.Certificate, caPrivKey *rsa.PrivateKey, etcdServers []*url.URL) ([]assetslib.Asset, error) {
	var assets []assetslib.Asset

	// Use the master CA to generate etcd assets.
	etcdCACert := caCert

	// Create an etcd client cert.
	etcdClientKey, etcdClientCert, err := r.newEtcdKeyAndCert(caCert, caPrivKey, "etcd-client", etcdServers)
	if err != nil {
		return nil, err
	}

	// Create an etcd peer cert (not consumed by self-hosted components).
	etcdPeerKey, etcdPeerCert, err := r.newEtcdKeyAndCert(caCert, caPrivKey, "etcd-peer", etcdServers)
	if err != nil {
		return nil, err
	}
	etcdServerKey, etcdServerCert, err := r.newEtcdKeyAndCert(caCert, caPrivKey, "etcd-server", etcdServers)
	if err != nil {
		return nil, err
	}

	const AssetPathEtcdPeerKey = ""
	assets = append(assets, []assetslib.Asset{
		{Name: AssetPathEtcdPeerCA, Data: tlsutil.EncodeCertificatePEM(etcdCACert)},
		{Name: AssetPathEtcdPeerKey, Data: tlsutil.EncodePrivateKeyPEM(etcdPeerKey)},
		{Name: AssetPathEtcdPeerCert, Data: tlsutil.EncodeCertificatePEM(etcdPeerCert)},
		{Name: AssetPathEtcdServerCA, Data: tlsutil.EncodeCertificatePEM(etcdCACert)},
		{Name: AssetPathEtcdServerKey, Data: tlsutil.EncodePrivateKeyPEM(etcdServerKey)},
		{Name: AssetPathEtcdServerCert, Data: tlsutil.EncodeCertificatePEM(etcdServerCert)},
	}...)

	assets = append(assets, []assetslib.Asset{
		{Name: AssetPathEtcdClientCA, Data: tlsutil.EncodeCertificatePEM(etcdCACert)},
		{Name: AssetPathEtcdClientKey, Data: tlsutil.EncodePrivateKeyPEM(etcdClientKey)},
		{Name: AssetPathEtcdClientCert, Data: tlsutil.EncodeCertificatePEM(etcdClientCert)},
	}...)

	return assets, nil
}

func (r *RenderOptions) newEtcdKeyAndCert(caCert *x509.Certificate, caPrivKey *rsa.PrivateKey, commonName string, etcdServers []*url.URL) (*rsa.PrivateKey, *x509.Certificate, error) {
	addrs := make([]string, len(etcdServers))
	for i := range etcdServers {
		addrs[i] = etcdServers[i].Hostname()
	}
	return r.newKeyAndCert(caCert, caPrivKey, commonName, addrs)
}
