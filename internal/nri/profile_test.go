package nri_test

import (
	api "github.com/containerd/nri/pkg/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nri "github.com/k8s-nono/nono-nri/internal/nri"
)

var _ = Describe("ResolveProfile", func() {
	cfg := &nri.Config{
		RuntimeClasses: []string{"nono-runc"},
		DefaultProfile: "default",
	}

	It("returns annotation value when present", func() {
		pod := &api.PodSandbox{
			Annotations: map[string]string{"nono.sh/profile": "strict"},
		}
		Expect(nri.ResolveProfile(pod, cfg)).To(Equal("strict"))
	})

	It("returns default profile when annotation absent", func() {
		pod := &api.PodSandbox{}
		Expect(nri.ResolveProfile(pod, cfg)).To(Equal("default"))
	})

	It("returns default profile when annotation is empty string", func() {
		pod := &api.PodSandbox{
			Annotations: map[string]string{"nono.sh/profile": ""},
		}
		Expect(nri.ResolveProfile(pod, cfg)).To(Equal("default"))
	})

	It("returns default profile when annotations map is nil", func() {
		pod := &api.PodSandbox{Annotations: nil}
		Expect(nri.ResolveProfile(pod, cfg)).To(Equal("default"))
	})
})
