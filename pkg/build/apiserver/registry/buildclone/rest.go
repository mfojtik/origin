package buildclone

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	buildv1 "github.com/openshift/api/build/v1"
	buildapi "github.com/openshift/origin/pkg/build/apis/build"
	"github.com/openshift/origin/pkg/build/generator"
)

// NewStorage creates a new storage object for build generation
func NewStorage(generator *generator.BuildGenerator) *CloneREST {
	return &CloneREST{generator: generator}
}

// CloneREST is a RESTStorage implementation for a BuildGenerator which supports only
// the Get operation (as the generator has no underlying storage object).
type CloneREST struct {
	generator *generator.BuildGenerator
}

var _ rest.Creater = &CloneREST{}

// New creates a new build clone request
func (s *CloneREST) New() runtime.Object {
	return &buildapi.BuildRequest{}
}

// Create instantiates a new build from an existing build
func (s *CloneREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, _ bool) (runtime.Object, error) {
	if err := rest.BeforeCreate(Strategy, ctx, obj); err != nil {
		return nil, err
	}
	if err := createValidation(obj); err != nil {
		return nil, err
	}
	internalRequest := obj.(*buildapi.BuildRequest)
	externalRequest := &buildv1.BuildRequest{}

	if err := legacyscheme.Scheme.Convert(internalRequest, externalRequest, nil); err != nil {
		return nil, err
	}
	externalClone, err := s.generator.Clone(ctx, externalRequest)
	if err != nil {
		return nil, err
	}

	result := &buildapi.BuildRequest{}
	err = legacyscheme.Scheme.Convert(externalClone, result, nil)
	return result, err
}
