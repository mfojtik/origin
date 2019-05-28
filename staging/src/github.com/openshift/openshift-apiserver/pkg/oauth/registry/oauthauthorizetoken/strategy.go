package oauthauthorizetoken

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	scopeauthorizer "github.com/openshift/openshift-apiserver/cmd/authorizer/scope"
	oauthapi "github.com/openshift/openshift-apiserver/pkg/oauth/apis/oauth"
	"github.com/openshift/openshift-apiserver/pkg/oauth/apis/oauth/validation"
	"github.com/openshift/openshift-apiserver/pkg/oauth/registry/oauthclient"
)

// strategy implements behavior for OAuthAuthorizeTokens
type strategy struct {
	runtime.ObjectTyper

	clientGetter oauthclient.Getter
}

var _ rest.RESTCreateStrategy = strategy{}
var _ rest.RESTUpdateStrategy = strategy{}
var _ rest.GarbageCollectionDeleteStrategy = strategy{}

func NewStrategy(clientGetter oauthclient.Getter) strategy {
	return strategy{ObjectTyper: legacyscheme.Scheme, clientGetter: clientGetter}
}

func (strategy) DefaultGarbageCollectionPolicy(ctx context.Context) rest.GarbageCollectionPolicy {
	return rest.Unsupported
}

func (strategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {}

// NamespaceScoped is false for OAuth objects
func (strategy) NamespaceScoped() bool {
	return false
}

func (strategy) GenerateName(base string) string {
	return base
}

func (strategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

// Canonicalize normalizes the object after validation.
func (strategy) Canonicalize(obj runtime.Object) {
}

// Validate validates a new token
func (s strategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	token := obj.(*oauthapi.OAuthAuthorizeToken)
	validationErrors := validation.ValidateAuthorizeToken(token)

	client, err := s.clientGetter.Get(token.ClientName, metav1.GetOptions{})
	if err != nil {
		return append(validationErrors, field.InternalError(field.NewPath("clientName"), err))
	}
	if err := scopeauthorizer.ValidateScopeRestrictions(client, token.Scopes...); err != nil {
		return append(validationErrors, field.InternalError(field.NewPath("clientName"), err))
	}

	return validationErrors
}

// ValidateUpdate validates an update
func (s strategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	oldToken := old.(*oauthapi.OAuthAuthorizeToken)
	newToken := obj.(*oauthapi.OAuthAuthorizeToken)
	return validation.ValidateAuthorizeTokenUpdate(newToken, oldToken)
}

// AllowCreateOnUpdate is false for OAuth objects
func (strategy) AllowCreateOnUpdate() bool {
	return false
}

func (strategy) AllowUnconditionalUpdate() bool {
	return false
}
