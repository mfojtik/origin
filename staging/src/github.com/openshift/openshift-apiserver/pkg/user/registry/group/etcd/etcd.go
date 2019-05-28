package etcd

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/kubernetes/pkg/printers"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"

	"github.com/openshift/api/user"
	userapi "github.com/openshift/openshift-apiserver/pkg/user/apis/user"
	printersinternal "github.com/openshift/openshift-apiserver/printers/internalversion"
	"github.com/openshift/origin/pkg/user/registry/group"
)

// REST implements a RESTStorage for groups against etcd
type REST struct {
	*registry.Store
}

// NewREST returns a RESTStorage object that will work against groups
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &registry.Store{
		NewFunc:                  func() runtime.Object { return &userapi.Group{} },
		NewListFunc:              func() runtime.Object { return &userapi.GroupList{} },
		DefaultQualifiedResource: user.Resource("groups"),

		TableConvertor: printerstorage.TableConvertor{TablePrinter: printers.NewTablePrinter().With(printersinternal.AddHandlers)},

		CreateStrategy: group.Strategy,
		UpdateStrategy: group.Strategy,
		DeleteStrategy: group.Strategy,
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	return &REST{store}, nil
}
