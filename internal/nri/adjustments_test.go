package nri_test

import (
	api "github.com/containerd/nri/pkg/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nri "github.com/k8s-nono/nono-nri/internal/nri"
)

var _ = Describe("BuildAdjustment", func() {
	DescribeTable("args prepend",
		func(originalArgs []string, profile string, expectedArgs []string) {
			ctr := &api.Container{Id: "ctr-x", Args: originalArgs}
			adj := nri.BuildAdjustment(ctr, profile, "/host/nono", false)
			Expect(adj.Args).To(Equal(expectedArgs))
		},
		Entry("with existing args",
			[]string{"myapp", "--port", "8080"}, "strict",
			[]string{"/nono/nono", "wrap", "--profile", "strict", "--", "myapp", "--port", "8080"},
		),
		Entry("with nil args",
			nil, "default",
			[]string{"/nono/nono", "wrap", "--profile", "default", "--"},
		),
		Entry("with empty args slice",
			[]string{}, "permissive",
			[]string{"/nono/nono", "wrap", "--profile", "permissive", "--"},
		),
	)

	Describe("readonly bind mount", func() {
		It("mounts the host directory to /nono so binary is accessible at /nono/nono", func() {
			ctr := &api.Container{Id: "ctr-3", Args: []string{"cmd"}}
			adj := nri.BuildAdjustment(ctr, "strict", "/usr/local/bin/nono", false)

			Expect(adj.Mounts).To(HaveLen(1))
			m := adj.Mounts[0]
			Expect(m.Source).To(Equal("/usr/local/bin"))
			Expect(m.Destination).To(Equal("/nono"))
			Expect(m.Type).To(Equal("bind"))
			Expect(m.Options).To(ContainElements("bind", "ro", "rprivate"))
		})

		It("skips bind-mount and injects NONO_PROFILE env when nonoInVMRootfs is true", func() {
			ctr := &api.Container{Id: "ctr-4", Args: []string{"cmd"}}
			adj := nri.BuildAdjustment(ctr, "strict", "", true)

			Expect(adj.Mounts).To(BeEmpty())
			Expect(adj.Env).To(HaveLen(1))
			Expect(adj.Env[0].Key).To(Equal("NONO_PROFILE"))
			Expect(adj.Env[0].Value).To(Equal("strict"))
		})
	})

	Describe("ContainerNonoPath constant", func() {
		It("equals /nono/nono", func() {
			Expect(nri.ContainerNonoPath).To(Equal("/nono/nono"))
		})
	})
})
