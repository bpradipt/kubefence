package nri

import (
	"path/filepath"

	api "github.com/containerd/nri/pkg/api"
)

// ContainerNonoPath is the fixed path inside the container where the nono
// binary is accessible after the host directory is bind-mounted.
const ContainerNonoPath = "/nono/nono"

// containerNonoDirPath is the directory bind-mounted into the container.
// Mounting the directory (not the file) ensures the destination path is
// created by the OCI runtime even when it doesn't pre-exist in the rootfs.
//
// Invariant: no container image used with nono-nri should define its own
// /nono directory for a different purpose. If a container image already has
// /nono, this bind-mount will silently shadow it. This is an accepted
// constraint of the current design.
const containerNonoDirPath = "/nono"

// BuildAdjustment constructs a ContainerAdjustment that wraps the container
// process with nono. It always prepends the nono wrapper to process.args.
//
// Standard delivery (vmRootfs=false): bind-mounts the directory containing
// hostBinPath into the container at /nono so the binary is accessible at
// ContainerNonoPath (/nono/nono).
//
// VM rootfs delivery (vmRootfs=true): the bind-mount is skipped (nono is
// pre-installed in the VM guest rootfs at /nono/nono). NONO_PROFILE is
// injected as an env var so shell wrapper scripts invoked via kubectl exec
// can pick up the correct profile without a separate config file.
//
// It is safe to call with a container that has nil or empty args.
func BuildAdjustment(ctr *api.Container, profile, hostBinPath string, vmRootfs bool) *api.ContainerAdjustment {
	prefix := []string{ContainerNonoPath, "wrap", "--profile", profile, "--"}
	orig := ctr.GetArgs()
	newArgs := make([]string, 0, len(prefix)+len(orig))
	newArgs = append(newArgs, prefix...)
	newArgs = append(newArgs, orig...)

	adj := &api.ContainerAdjustment{}
	adj.SetArgs(newArgs)

	if vmRootfs {
		adj.AddEnv("NONO_PROFILE", profile)
	} else {
		adj.AddMount(&api.Mount{
			Source:      filepath.Dir(hostBinPath),
			Destination: containerNonoDirPath,
			Type:        "bind",
			Options:     []string{"bind", "ro", "rprivate"},
		})
	}
	return adj
}
