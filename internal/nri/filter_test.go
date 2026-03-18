package nri_test

import (
	api "github.com/containerd/nri/pkg/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nri "github.com/k8s-nono/nono-nri/internal/nri"
)

var _ = Describe("ShouldSandbox", func() {
	cfg := &nri.Config{
		RuntimeClasses: []string{"nono-runc", "nono-kata"},
		DefaultProfile: "default",
	}

	DescribeTable("filter decision",
		func(handler string, expectSandbox bool) {
			pod := &api.PodSandbox{RuntimeHandler: handler}
			Expect(nri.ShouldSandbox(pod, cfg)).To(Equal(expectSandbox))
		},
		Entry("matching handler nono-runc", "nono-runc", true),
		Entry("matching handler nono-kata", "nono-kata", true),
		Entry("non-matching handler runc", "runc", false),
		Entry("empty handler", "", false),
		Entry("similar but wrong handler nono-runc-extra", "nono-runc-extra", false),
	)
})

var _ = Describe("SkipReason", func() {
	It("returns no_runtime_handler for empty handler", func() {
		pod := &api.PodSandbox{RuntimeHandler: ""}
		Expect(nri.SkipReason(pod)).To(Equal("no_runtime_handler"))
	})

	It("returns runtime_class_not_matched for non-matching handler", func() {
		pod := &api.PodSandbox{RuntimeHandler: "runc"}
		Expect(nri.SkipReason(pod)).To(Equal("runtime_class_not_matched"))
	})
})
