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
const containerNonoDirPath = "/nono"

// BuildAdjustment constructs a ContainerAdjustment that:
//  1. Prepends the nono wrapper command prefix to the container's existing args
//     via SetArgs: [ContainerNonoPath, "wrap", "--profile", profile, "--", <original args...>]
//  2. Adds a readonly bind mount of the directory containing hostBinPath to
//     containerNonoDirPath, making the binary accessible at ContainerNonoPath.
//
// It is safe to call with a container that has nil or empty args.
func BuildAdjustment(ctr *api.Container, profile, hostBinPath string) *api.ContainerAdjustment {
	prefix := []string{ContainerNonoPath, "wrap", "--profile", profile, "--"}
	orig := ctr.GetArgs()
	newArgs := make([]string, 0, len(prefix)+len(orig))
	newArgs = append(newArgs, prefix...)
	newArgs = append(newArgs, orig...)

	adj := &api.ContainerAdjustment{}
	adj.SetArgs(newArgs)
	adj.AddMount(&api.Mount{
		Source:      filepath.Dir(hostBinPath),
		Destination: containerNonoDirPath,
		Type:        "bind",
		Options:     []string{"bind", "ro", "rprivate"},
	})
	return adj
}
