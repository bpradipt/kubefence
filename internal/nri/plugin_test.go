package nri_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"

	api "github.com/containerd/nri/pkg/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	nri "github.com/k8s-nono/nono-nri/internal/nri"
)

// logEntry is used to parse structured JSON log output from the plugin.
type logEntry struct {
	Msg            string `json:"msg"`
	Decision       string `json:"decision"`
	ContainerID    string `json:"container_id"`
	Namespace      string `json:"namespace"`
	Pod            string `json:"pod"`
	Profile        string `json:"profile"`
	RuntimeHandler string `json:"runtime_handler"`
	Reason         string `json:"reason"`
}

// newBufLogger creates a JSON slog.Logger that writes to the returned buffer.
func newBufLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

var _ = Describe("Plugin", func() {
	Describe("CreateContainer", func() {
		It("logs skip with all required fields for non-matching container", func() {
			cfg := &nri.Config{
				RuntimeClasses: []string{"nono-runc"},
				DefaultProfile: "default",
			}
			buf := &bytes.Buffer{}
			p := nri.NewPlugin(cfg, newBufLogger(buf))

			pod := &api.PodSandbox{
				RuntimeHandler: "runc",
				Namespace:      "prod",
				Name:           "web-abc",
				Annotations:    map[string]string{},
			}
			ctr := &api.Container{Id: "ctr-123"}

			adj, updates, err := p.CreateContainer(context.Background(), pod, ctr)
			Expect(err).To(BeNil())
			Expect(adj).To(BeNil())
			Expect(updates).To(BeNil())

			var entry logEntry
			Expect(json.Unmarshal(buf.Bytes(), &entry)).To(Succeed())
			Expect(entry.Msg).To(Equal("skip"))
			Expect(entry.Decision).To(Equal("skip"))
			Expect(entry.ContainerID).To(Equal("ctr-123"))
			Expect(entry.Namespace).To(Equal("prod"))
			Expect(entry.Pod).To(Equal("web-abc"))
			Expect(entry.RuntimeHandler).To(Equal("runc"))
			Expect(entry.Reason).NotTo(BeEmpty())
		})

		It("logs injection-pending with all required fields for matching container", func() {
			cfg := &nri.Config{
				RuntimeClasses: []string{"nono-runc"},
				DefaultProfile: "default",
			}
			buf := &bytes.Buffer{}
			p := nri.NewPlugin(cfg, newBufLogger(buf))

			pod := &api.PodSandbox{
				RuntimeHandler: "nono-runc",
				Namespace:      "prod",
				Name:           "app-xyz",
				Annotations:    map[string]string{"nono.sh/profile": "strict"},
			}
			ctr := &api.Container{Id: "ctr-456"}

			adj, updates, err := p.CreateContainer(context.Background(), pod, ctr)
			Expect(err).To(BeNil())
			Expect(adj).To(BeNil())
			Expect(updates).To(BeNil())

			var entry logEntry
			Expect(json.Unmarshal(buf.Bytes(), &entry)).To(Succeed())
			Expect(entry.Msg).To(Equal("injection-pending"))
			Expect(entry.Decision).To(Equal("inject"))
			Expect(entry.ContainerID).To(Equal("ctr-456"))
			Expect(entry.Namespace).To(Equal("prod"))
			Expect(entry.Pod).To(Equal("app-xyz"))
			Expect(entry.Profile).To(Equal("strict"))
			Expect(entry.RuntimeHandler).To(Equal("nono-runc"))
		})

		It("returns nil ContainerAdjustment (no-op)", func() {
			cfg := &nri.Config{
				RuntimeClasses: []string{"nono-runc"},
				DefaultProfile: "default",
			}
			buf := &bytes.Buffer{}
			p := nri.NewPlugin(cfg, newBufLogger(buf))

			pod := &api.PodSandbox{
				RuntimeHandler: "nono-runc",
				Namespace:      "default",
				Name:           "test-pod",
				Annotations:    map[string]string{},
			}
			ctr := &api.Container{Id: "ctr-789"}

			adj, _, _ := p.CreateContainer(context.Background(), pod, ctr)
			Expect(adj).To(BeNil())
		})

		It("uses default profile when annotation absent", func() {
			cfg := &nri.Config{
				RuntimeClasses: []string{"nono-runc"},
				DefaultProfile: "fallback",
			}
			buf := &bytes.Buffer{}
			p := nri.NewPlugin(cfg, newBufLogger(buf))

			pod := &api.PodSandbox{
				RuntimeHandler: "nono-runc",
				Namespace:      "default",
				Name:           "test-pod",
				Annotations:    map[string]string{},
			}
			ctr := &api.Container{Id: "ctr-000"}

			_, _, err := p.CreateContainer(context.Background(), pod, ctr)
			Expect(err).To(BeNil())

			var entry logEntry
			Expect(json.Unmarshal(buf.Bytes(), &entry)).To(Succeed())
			Expect(entry.Profile).To(Equal("fallback"))
		})
	})

	Describe("RemoveContainer", func() {
		It("returns nil updates and nil error", func() {
			cfg := &nri.Config{
				RuntimeClasses: []string{"nono-runc"},
				DefaultProfile: "default",
			}
			buf := &bytes.Buffer{}
			p := nri.NewPlugin(cfg, newBufLogger(buf))

			pod := &api.PodSandbox{
				Name:      "test-pod",
				Namespace: "default",
			}
			ctr := &api.Container{Id: "ctr-rem"}

			updates, err := p.RemoveContainer(context.Background(), pod, ctr)
			Expect(err).To(BeNil())
			Expect(updates).To(BeNil())
		})
	})
})
