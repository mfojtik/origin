package etcd

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/kubernetes/pkg/printers"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"

	"github.com/openshift/api/user"
	userapi "github.com/openshift/openshift-apiserver/pkg/user/apis/user"
	"github.com/openshift/openshift-apiserver/pkg/user/registry/identity"
	printersinternal "github.com/openshift/openshift-apiserver/printers/internalversion"
)

// REST implements a RESTStorage for identites against etcd
type REST struct {
	*registry.Store
}

var _ rest.StandardStorage = &REST{}

// NewREST returns a RESTStorage object that will work against identites
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &userapi.Identity{} },
		NewListFunc:              func() runtime.Object { return &userapi.IdentityList{} },
		DefaultQualifiedResource: user.Resource("identities"),

		TableConvertor: printerstorage.TableConvertor{TablePrinter: printers.NewTablePrinter().With(printersinternal.AddHandlers)},

		CreateStrategy: identity.Strategy,
		UpdateStrategy: identity.Strategy,
		DeleteStrategy: identity.Strategy,
	}

	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    storage.AttrFunc(storage.DefaultNamespaceScopedAttr).WithFieldMutation(userapi.IdentityFieldSelector),
	}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil
}
