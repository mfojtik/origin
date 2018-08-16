package build

import (
	kapi "k8s.io/kubernetes/pkg/apis/core"

	"github.com/openshift/origin/pkg/api/apihelpers"
)

// NOTE: These helpers are used by apiserver only as the apiserver use the internal types.
//       These were copied from pkg/build/util and any change to original helpers should be reflected here as well.

const (
	// buildPodSuffix is the suffix used to append to a build pod name given a build name
	buildPodSuffix = "build"
)

// DEPRECATED: Reserved for apiserver, do not use outside of it
func IsInternalBuildComplete(b *Build) bool {
	return IsInternalTerminalPhase(b.Status.Phase)
}

// DEPRECATED: Reserved for apiserver, do not use outside of it
func IsInternalTerminalPhase(p BuildPhase) bool {
	switch p {
	case BuildPhaseNew,
		BuildPhasePending,
		BuildPhaseRunning:
		return false
	}
	return true
}

// GetBuildPodName returns name of the build pod.
// DEPRECATED: Reserved for apiserver, do not use outside of it
func GetInternalBuildPodName(build *Build) string {
	return apihelpers.GetPodName(build.Name, buildPodSuffix)
}

// DEPRECATED: Reserved for apiserver, do not use outside of it
func InternalStrategyType(strategy BuildStrategy) string {
	switch {
	case strategy.DockerStrategy != nil:
		return "Docker"
	case strategy.CustomStrategy != nil:
		return "Custom"
	case strategy.SourceStrategy != nil:
		return "Source"
	case strategy.JenkinsPipelineStrategy != nil:
		return "JenkinsPipeline"
	}
	return ""
}

// GetInternalInputReference returns the From ObjectReference associated with the
// BuildStrategy.
// DEPRECATED: Reserved for apiserver, do not use outside of it
func GetInternalInputReference(strategy BuildStrategy) *kapi.ObjectReference {
	switch {
	case strategy.SourceStrategy != nil:
		return &strategy.SourceStrategy.From
	case strategy.DockerStrategy != nil:
		return strategy.DockerStrategy.From
	case strategy.CustomStrategy != nil:
		return &strategy.CustomStrategy.From
	default:
		return nil
	}
}

// GetInternalBuildEnv gets the build strategy environment
// DEPRECATED: Reserved for apiserver, do not use outside of it
func GetInternalBuildEnv(build *Build) []kapi.EnvVar {
	switch {
	case build.Spec.Strategy.SourceStrategy != nil:
		return build.Spec.Strategy.SourceStrategy.Env
	case build.Spec.Strategy.DockerStrategy != nil:
		return build.Spec.Strategy.DockerStrategy.Env
	case build.Spec.Strategy.CustomStrategy != nil:
		return build.Spec.Strategy.CustomStrategy.Env
	case build.Spec.Strategy.JenkinsPipelineStrategy != nil:
		return build.Spec.Strategy.JenkinsPipelineStrategy.Env
	default:
		return nil
	}
}

// SetInternalBuildEnv replaces the current build environment
// DEPRECATED: Reserved for apiserver, do not use outside of it
func SetInternalBuildEnv(build *Build, env []kapi.EnvVar) {
	var oldEnv *[]kapi.EnvVar

	switch {
	case build.Spec.Strategy.SourceStrategy != nil:
		oldEnv = &build.Spec.Strategy.SourceStrategy.Env
	case build.Spec.Strategy.DockerStrategy != nil:
		oldEnv = &build.Spec.Strategy.DockerStrategy.Env
	case build.Spec.Strategy.CustomStrategy != nil:
		oldEnv = &build.Spec.Strategy.CustomStrategy.Env
	case build.Spec.Strategy.JenkinsPipelineStrategy != nil:
		oldEnv = &build.Spec.Strategy.JenkinsPipelineStrategy.Env
	default:
		return
	}
	*oldEnv = env
}

// UpdateInternalBuildEnv updates the strategy environment
// This will replace the existing variable definitions with provided env
// DEPRECATED: Reserved for apiserver, do not use outside of it
func UpdateInternalBuildEnv(build *Build, env []kapi.EnvVar) {
	buildEnv := GetInternalBuildEnv(build)

	newEnv := []kapi.EnvVar{}
	for _, e := range buildEnv {
		exists := false
		for _, n := range env {
			if e.Name == n.Name {
				exists = true
				break
			}
		}
		if !exists {
			newEnv = append(newEnv, e)
		}
	}
	newEnv = append(newEnv, env...)
	SetInternalBuildEnv(build, newEnv)
}

// InternalBuildSliceByCreationTimestamp implements sort.Interface for []Build
// based on the CreationTimestamp field.
// DEPRECATED: Reserved for apiserver, do not use outside of it
// +k8s:deepcopy-gen=false
type InternalBuildSliceByCreationTimestamp []Build

func (b InternalBuildSliceByCreationTimestamp) Len() int {
	return len(b)
}

func (b InternalBuildSliceByCreationTimestamp) Less(i, j int) bool {
	return b[i].CreationTimestamp.Before(&b[j].CreationTimestamp)
}

func (b InternalBuildSliceByCreationTimestamp) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
