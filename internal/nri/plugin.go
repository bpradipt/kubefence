package nri

import (
	"context"
	"log/slog"

	api "github.com/containerd/nri/pkg/api"
)

// Plugin implements the NRI plugin interface for nono-nri.
// It intercepts container creation events, decides whether to sandbox the
// container based on the pod's RuntimeHandler, and logs structured decisions
// with all CORE-04 required fields.
type Plugin struct {
	Config *Config
	Log    *slog.Logger
}

// NewPlugin constructs a Plugin with the given config and logger.
func NewPlugin(cfg *Config, logger *slog.Logger) *Plugin {
	return &Plugin{Config: cfg, Log: logger}
}

// CreateContainer is called by the NRI runtime before each container is created.
// It resolves the nono profile for the pod, checks whether the container should
// be sandboxed, and logs the resulting decision with all CORE-04 fields.
// Phase 1 is a no-op: it always returns a nil ContainerAdjustment.
func (p *Plugin) CreateContainer(
	ctx context.Context,
	pod *api.PodSandbox,
	ctr *api.Container,
) (*api.ContainerAdjustment, []*api.ContainerUpdate, error) {
	handler := pod.GetRuntimeHandler()
	namespace := pod.GetNamespace()
	podName := pod.GetName()
	ctrID := ctr.GetId()
	profile := ResolveProfile(pod, p.Config)

	if !ShouldSandbox(pod, p.Config) {
		p.Log.Info("skip",
			"decision", "skip",
			"reason", SkipReason(pod, p.Config),
			"container_id", ctrID,
			"namespace", namespace,
			"pod", podName,
			"profile", profile,
			"runtime_handler", handler,
		)
		return nil, nil, nil
	}

	p.Log.Info("injection-pending",
		"decision", "inject",
		"container_id", ctrID,
		"namespace", namespace,
		"pod", podName,
		"profile", profile,
		"runtime_handler", handler,
	)
	return nil, nil, nil
}

// RemoveContainer is called by the NRI runtime after a container is removed.
// It logs the removal event at debug level and returns no updates.
func (p *Plugin) RemoveContainer(
	ctx context.Context,
	pod *api.PodSandbox,
	ctr *api.Container,
) ([]*api.ContainerUpdate, error) {
	p.Log.Debug("container-removed",
		"container_id", ctr.GetId(),
		"pod", pod.GetName(),
		"namespace", pod.GetNamespace(),
	)
	return nil, nil
}
