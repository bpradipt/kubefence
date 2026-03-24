package nri

import (
	"regexp"

	api "github.com/containerd/nri/pkg/api"
)

// ProfileAnnotationKey is the pod annotation key that specifies the nono profile.
const ProfileAnnotationKey = "nono.sh/profile"

// validProfileRe restricts profile names to safe identifiers that must start
// with an alphanumeric character, followed by up to 63 alphanumeric, hyphen,
// or underscore characters. The leading-alphanumeric requirement prevents
// values like "--allow-all" from being treated as CLI flags by the nono binary.
var validProfileRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,63}$`)

// ResolveProfile returns the nono profile for the given pod.
// It returns the annotation value from ProfileAnnotationKey if present and it
// matches [a-zA-Z0-9_-]{1,64}, otherwise it falls back to cfg.DefaultProfile.
func ResolveProfile(pod *api.PodSandbox, cfg *Config) string {
	annotations := pod.GetAnnotations()
	if annotations != nil {
		if profile, ok := annotations[ProfileAnnotationKey]; ok && validProfileRe.MatchString(profile) {
			return profile
		}
	}
	return cfg.DefaultProfile
}
