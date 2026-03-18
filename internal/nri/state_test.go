package nri_test

import (
	"encoding/json"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nri "github.com/k8s-nono/nono-nri/internal/nri"
)

var _ = Describe("State", func() {
	var tmpDir string

	BeforeEach(func() {
		var err error
		tmpDir, err = os.MkdirTemp("", "nono-state-test-*")
		Expect(err).NotTo(HaveOccurred())
		nri.SetStateBaseDir(tmpDir)
	})

	AfterEach(func() {
		nri.ResetStateBaseDir()
		os.RemoveAll(tmpDir)
	})

	Describe("WriteMetadata", func() {
		It("creates the container state dir and metadata.json", func() {
			err := nri.WriteMetadata("pod-uid-1", "ctr-1", "web", "prod", "strict")
			Expect(err).NotTo(HaveOccurred())

			expectedDir := filepath.Join(tmpDir, "pod-uid-1", "ctr-1")
			info, statErr := os.Stat(expectedDir)
			Expect(statErr).NotTo(HaveOccurred())
			Expect(info.IsDir()).To(BeTrue())

			metaPath := filepath.Join(expectedDir, "metadata.json")
			_, statErr = os.Stat(metaPath)
			Expect(statErr).NotTo(HaveOccurred())
		})

		It("writes valid JSON with all required fields", func() {
			err := nri.WriteMetadata("pod-uid-2", "ctr-2", "web", "prod", "strict")
			Expect(err).NotTo(HaveOccurred())

			metaPath := filepath.Join(tmpDir, "pod-uid-2", "ctr-2", "metadata.json")
			data, readErr := os.ReadFile(metaPath)
			Expect(readErr).NotTo(HaveOccurred())

			var meta nri.ContainerMetadata
			Expect(json.Unmarshal(data, &meta)).To(Succeed())

			Expect(meta.ContainerID).To(Equal("ctr-2"))
			Expect(meta.Pod).To(Equal("web"))
			Expect(meta.Namespace).To(Equal("prod"))
			Expect(meta.Profile).To(Equal("strict"))
			Expect(meta.Timestamp).NotTo(BeEmpty())
		})

		It("rejects podUID containing a path separator", func() {
			err := nri.WriteMetadata("../../etc", "ctr-1", "web", "prod", "strict")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid path component"))
		})

		It("rejects containerID containing a path separator", func() {
			err := nri.WriteMetadata("pod-uid-1", "../escape", "web", "prod", "strict")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid path component"))
		})

		It("rejects containerID that is exactly ..", func() {
			err := nri.WriteMetadata("pod-uid-1", "..", "web", "prod", "strict")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid path component"))
		})
	})

	Describe("RemoveMetadata", func() {
		It("removes the container directory", func() {
			Expect(nri.WriteMetadata("pod-uid-3", "ctr-3", "web", "prod", "strict")).To(Succeed())

			err := nri.RemoveMetadata("pod-uid-3", "ctr-3")
			Expect(err).NotTo(HaveOccurred())

			ctrDir := filepath.Join(tmpDir, "pod-uid-3", "ctr-3")
			_, statErr := os.Stat(ctrDir)
			Expect(os.IsNotExist(statErr)).To(BeTrue())
		})

		It("removes the pod parent dir when it becomes empty", func() {
			Expect(nri.WriteMetadata("pod-uid-4", "ctr-4", "web", "prod", "strict")).To(Succeed())

			Expect(nri.RemoveMetadata("pod-uid-4", "ctr-4")).To(Succeed())

			podDir := filepath.Join(tmpDir, "pod-uid-4")
			_, statErr := os.Stat(podDir)
			Expect(os.IsNotExist(statErr)).To(BeTrue())
		})

		It("does NOT remove the pod dir when other container dirs still exist", func() {
			Expect(nri.WriteMetadata("pod-uid-5", "ctr-5a", "web", "prod", "strict")).To(Succeed())
			Expect(nri.WriteMetadata("pod-uid-5", "ctr-5b", "web", "prod", "strict")).To(Succeed())

			Expect(nri.RemoveMetadata("pod-uid-5", "ctr-5a")).To(Succeed())

			podDir := filepath.Join(tmpDir, "pod-uid-5")
			info, statErr := os.Stat(podDir)
			Expect(statErr).NotTo(HaveOccurred())
			Expect(info.IsDir()).To(BeTrue())
		})

		It("returns nil error for non-existent paths", func() {
			err := nri.RemoveMetadata("nonexistent-pod", "nonexistent-ctr")
			Expect(err).NotTo(HaveOccurred())
		})

		It("rejects podUID containing a path separator", func() {
			err := nri.RemoveMetadata("../../etc", "ctr-1")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid path component"))
		})
	})
})
