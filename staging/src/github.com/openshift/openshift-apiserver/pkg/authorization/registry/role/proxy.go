package role

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/printers"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"

	authorizationapi "github.com/openshift/openshift-apiserver/pkg/authorization/apis/authorization"
	"github.com/openshift/openshift-apiserver/pkg/authorization/apiserver/registry/util"
	authclient "github.com/openshift/origin/pkg/client/impersonatingclient"
	printersinternal "github.com/openshift/origin/pkg/printers/internalversion"
	utilregistry "github.com/openshift/origin/pkg/util/registry"
)

type REST struct {
	privilegedClient restclient.Interface
	rest.TableConvertor
}

var _ rest.Lister = &REST{}
var _ rest.Getter = &REST{}
var _ rest.CreaterUpdater = &REST{}
var _ rest.GracefulDeleter = &REST{}
var _ rest.Scoper = &REST{}

func NewREST(client restclient.Interface) utilregistry.NoWatchStorage {
	return utilregistry.WrapNoWatchStorageError(&REST{
		privilegedClient: client,
		TableConvertor:   printerstorage.TableConvertor{TablePrinter: printers.NewTablePrinter().With(printersinternal.AddHandlers)},
	})
}

func (s *REST) New() runtime.Object {
	return &authorizationapi.Role{}
}
func (s *REST) NewList() runtime.Object {
	return &authorizationapi.RoleList{}
}

func (s *REST) NamespaceScoped() bool {
	return true
}

func (s *REST) List(ctx context.Context, options *metainternal.ListOptions) (runtime.Object, error) {
	client, err := s.getImpersonatingClient(ctx)
	if err != nil {
		return nil, err
	}

	optv1 := metav1.ListOptions{}
	if err := metainternal.Convert_internalversion_ListOptions_To_v1_ListOptions(options, &optv1, nil); err != nil {
		return nil, err
	}

	roles, err := client.List(optv1)
	if err != nil {
		return nil, err
	}

	ret := &authorizationapi.RoleList{ListMeta: roles.ListMeta}
	for _, curr := range roles.Items {
		role, err := util.RoleFromRBAC(&curr)
		if err != nil {
			return nil, err
		}
		ret.Items = append(ret.Items, *role)
	}
	return ret, nil
}

func (s *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	client, err := s.getImpersonatingClient(ctx)
	if err != nil {
		return nil, err
	}

	ret, err := client.Get(name, *options)
	if err != nil {
		return nil, err
	}

	role, err := util.RoleFromRBAC(ret)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *REST) Delete(ctx context.Context, name string, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	client, err := s.getImpersonatingClient(ctx)
	if err != nil {
		return nil, false, err
	}

	if err := client.Delete(name, options); err != nil {
		return nil, false, err
	}

	return &metav1.Status{Status: metav1.StatusSuccess}, true, nil
}

func (s *REST) Create(ctx context.Context, obj runtime.Object, _ rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	client, err := s.getImpersonatingClient(ctx)
	if err != nil {
		return nil, err
	}

	convertedObj, err := util.RoleToRBAC(obj.(*authorizationapi.Role))
	if err != nil {
		return nil, err
	}

	ret, err := client.Create(convertedObj)
	if err != nil {
		return nil, err
	}

	role, err := util.RoleFromRBAC(ret)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *REST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, _ rest.ValidateObjectFunc, _ rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	client, err := s.getImpersonatingClient(ctx)
	if err != nil {
		return nil, false, err
	}

	old, err := client.Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, false, err
	}

	oldRole, err := util.RoleFromRBAC(old)
	if err != nil {
		return nil, false, err
	}

	obj, err := objInfo.UpdatedObject(ctx, oldRole)
	if err != nil {
		return nil, false, err
	}

	updatedRole, err := util.RoleToRBAC(obj.(*authorizationapi.Role))
	if err != nil {
		return nil, false, err
	}

	ret, err := client.Update(updatedRole)
	if err != nil {
		return nil, false, err
	}

	role, err := util.RoleFromRBAC(ret)
	if err != nil {
		return nil, false, err
	}
	return role, false, err
}

func (s *REST) getImpersonatingClient(ctx context.Context) (rbacv1.RoleInterface, error) {
	namespace, ok := apirequest.NamespaceFrom(ctx)
	if !ok {
		return nil, apierrors.NewBadRequest("namespace parameter required")
	}
	rbacClient, err := authclient.NewImpersonatingRBACFromContext(ctx, s.privilegedClient)
	if err != nil {
		return nil, err
	}
	return rbacClient.Roles(namespace), nil
}
