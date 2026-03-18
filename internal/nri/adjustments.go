package nri

import (
	api "github.com/containerd/nri/pkg/api"
)

// ContainerNonoPath is the fixed path inside the container where the nono
// binary is bind-mounted from the host.
const ContainerNonoPath = "/nono/nono"

// BuildAdjustment constructs a ContainerAdjustment that:
//  1. Prepends the nono wrapper command prefix to the container's existing args
//     via SetArgs: [ContainerNonoPath, "wrap", "--profile", profile, "--", <original args...>]
//  2. Adds a readonly bind mount from hostBinPath on the host to ContainerNonoPath
//     inside the container.
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
		Source:      hostBinPath,
		Destination: ContainerNonoPath,
		Type:        "bind",
		Options:     []string{"bind", "ro", "rprivate"},
	})
	return adj
}
