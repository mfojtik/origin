package build

import (
	kapi "k8s.io/kubernetes/pkg/apis/core"

	"github.com/openshift/origin/pkg/api/apihelpers"
)

const (
	// buildPodSuffix is the suffix used to append to a build pod name given a build name
	buildPodSuffix = "build"
)

// DEPRECATED: Use external version
func IsInternalBuildComplete(b *Build) bool {
	return IsInternalTerminalPhase(b.Status.Phase)
}

// DEPRECATED: Use external version
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
// DEPRECATED: Use external version
func GetInternalBuildPodName(build *Build) string {
	return apihelpers.GetPodName(build.Name, buildPodSuffix)
}

// DEPRECATED: Use external version
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
// DEPRECATED: Use external version
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

// InternalBuildSliceByCreationTimestamp implements sort.Interface for []Build
// based on the CreationTimestamp field.
// DEPRECATED:
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
