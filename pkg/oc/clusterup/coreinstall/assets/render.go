package assets

import (
	"net/url"

	assetslib "github.com/openshift/library-go/pkg/assets"
	"github.com/openshift/origin/pkg/oc/clusterup/coreinstall/assets/tlsutil"
)

type RenderOptions struct {
	AltNames                *tlsutil.AltNames
	EtcdServerURLs          []*url.URL
	CommonName              string
	OrganizationalUnitNames []string
}

func NewRenderer() *RenderOptions {
	return &RenderOptions{
		CommonName:              "kube-ca",
		OrganizationalUnitNames: []string{"bootkube"},
	}
}

func (r *RenderOptions) Render() (*assetslib.Assets, error) {
	var (
		err error
	)
	result := assetslib.Assets{}

	cAPrivKey, cACert, err := r.newCACert()
	if err != nil {
		return nil, err
	}

	if files, err := r.newTLSAssets(cACert, cAPrivKey, *r.AltNames); err != nil {
		return nil, err
	} else {
		result = append(result, files...)
	}

	if files, err := r.newEtcdTLSAssets(cACert, cAPrivKey, r.EtcdServerURLs); err != nil {
		return nil, err
	} else {
		result = append(result, files...)
	}

	return &result, nil
}
