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

	DescribeTable("rejects invalid annotation values and falls back to default",
		func(annotation string) {
			pod := &api.PodSandbox{
				Annotations: map[string]string{"nono.sh/profile": annotation},
			}
			Expect(nri.ResolveProfile(pod, cfg)).To(Equal("default"))
		},
		Entry("path separator", "../../etc/passwd"),
		Entry("forward slash", "a/b"),
		Entry("backslash", `a\b`),
		Entry("space", "my profile"),
		Entry("CLI flag injection", "--allow-all"),
		Entry("shell subshell", "$(evil)"),
		Entry("backtick", "`evil`"),
		Entry("semicolon", "foo;bar"),
		Entry("null byte", "foo\x00bar"),
		Entry("too long (65 chars)", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1"),
	)

	DescribeTable("accepts valid profile names",
		func(annotation string) {
			pod := &api.PodSandbox{
				Annotations: map[string]string{"nono.sh/profile": annotation},
			}
			Expect(nri.ResolveProfile(pod, cfg)).To(Equal(annotation))
		},
		Entry("lowercase letters", "strict"),
		Entry("uppercase letters", "Strict"),
		Entry("numbers", "profile1"),
		Entry("underscores", "my_profile"),
		Entry("hyphens", "my-profile"),
		Entry("mixed", "My-Profile_v2"),
		Entry("exactly 64 chars", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
	)
})
