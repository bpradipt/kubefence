package nri

import (
	api "github.com/containerd/nri/pkg/api"
)

// ProfileAnnotationKey is the pod annotation key that specifies the nono profile.
const ProfileAnnotationKey = "nono.sh/profile"

// ResolveProfile returns the nono profile for the given pod.
// It returns the annotation value from ProfileAnnotationKey if present and non-empty,
// otherwise it falls back to cfg.DefaultProfile.
func ResolveProfile(pod *api.PodSandbox, cfg *Config) string {
	annotations := pod.GetAnnotations()
	if annotations != nil {
		if profile, ok := annotations[ProfileAnnotationKey]; ok && profile != "" {
			return profile
		}
	}
	return cfg.DefaultProfile
}
