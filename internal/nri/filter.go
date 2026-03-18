package nri

import (
	api "github.com/containerd/nri/pkg/api"
)

// ShouldSandbox returns true if the pod should be sandboxed by nono-nri.
// A pod is sandboxed when its RuntimeHandler matches one of the configured RuntimeClasses.
// Pure function — no external I/O, easily unit-testable.
func ShouldSandbox(pod *api.PodSandbox, cfg *Config) bool {
	handler := pod.GetRuntimeHandler()
	for _, rc := range cfg.RuntimeClasses {
		if rc == handler {
			return true
		}
	}
	return false
}

// SkipReason returns a human-readable reason why ShouldSandbox returned false.
func SkipReason(pod *api.PodSandbox) string {
	if pod.GetRuntimeHandler() == "" {
		return "no_runtime_handler"
	}
	return "runtime_class_not_matched"
}
