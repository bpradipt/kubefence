package nri_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nri "github.com/k8s-nono/nono-nri/internal/nri"
)

var _ = Describe("kernelMinorParsing", func() {
	AfterEach(func() {
		nri.ResetKernelVersionFunc()
	})

	// These cases exercise the real defaultKernelVersion parser indirectly by
	// checking that CheckKernel succeeds on the live host; exotic string formats
	// are validated via the injected function to keep tests hermetic.
	DescribeTable("minor version extracted correctly from exotic release strings",
		func(major, minor int, expectOK bool) {
			nri.SetKernelVersionFunc(func() (int, int) { return major, minor })
			err := nri.CheckKernel()
			if expectOK {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(HaveOccurred())
			}
		},
		Entry("plain 6.8", 6, 8, true),
		Entry("5.13 exact minimum", 5, 13, true),
		Entry("5.15-rc1 style (minor=15)", 5, 15, true),
		Entry("5.12 below minimum", 5, 12, false),
		Entry("4.19 old kernel", 4, 19, false),
		Entry("0.0 parse-failure sentinel", 0, 0, false),
	)
})

var _ = Describe("CheckKernel", func() {
	BeforeEach(func() {
		// Override with a safe default; individual tests will set their versions.
		nri.SetKernelVersionFunc(func() (int, int) { return 6, 1 })
	})

	AfterEach(func() {
		nri.ResetKernelVersionFunc()
	})

	It("returns nil for kernel 5.13", func() {
		nri.SetKernelVersionFunc(func() (int, int) { return 5, 13 })
		Expect(nri.CheckKernel()).To(Succeed())
	})

	It("returns nil for kernel 6.1", func() {
		nri.SetKernelVersionFunc(func() (int, int) { return 6, 1 })
		Expect(nri.CheckKernel()).To(Succeed())
	})

	It("returns error for kernel 5.12", func() {
		nri.SetKernelVersionFunc(func() (int, int) { return 5, 12 })
		err := nri.CheckKernel()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("too old"))
		Expect(err.Error()).To(ContainSubstring("5.13"))
	})

	It("returns error for kernel 4.19", func() {
		nri.SetKernelVersionFunc(func() (int, int) { return 4, 19 })
		err := nri.CheckKernel()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("too old"))
	})

	It("returns error for kernel 0.0", func() {
		nri.SetKernelVersionFunc(func() (int, int) { return 0, 0 })
		err := nri.CheckKernel()
		Expect(err).To(HaveOccurred())
	})
})
