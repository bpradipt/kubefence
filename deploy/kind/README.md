# Kind Local Test Cluster

## Prerequisites

- Docker (running)
- [Kind](https://kind.sigs.k8s.io/) v0.27+
- kubectl
- Go 1.24+ (for building the plugin)
- `nono` binary placed at repo root as `./nono` (download from nono releases)

## Quick Start

One-command deployment:
```bash
bash deploy/kind/deploy.sh
```

This will:
1. Create a Kind cluster named `nono-test` with NRI enabled
2. Build the `nono-nri:latest` Docker image
3. Load the image into the Kind cluster
4. Copy the plugin TOML config to the Kind node
5. Apply RuntimeClass and DaemonSet manifests
6. Wait for the DaemonSet to become ready

## Verify Deployment

Check DaemonSet status:
```bash
kubectl get daemonset -n kube-system nono-nri
kubectl logs -n kube-system -l app=nono-nri
```

Test nono injection with a sample pod:
```bash
kubectl apply -f deploy/test-pod.yaml
kubectl wait --for=condition=ready pod/nono-test --timeout=60s
kubectl exec nono-test -- cat /proc/1/cmdline | tr '\0' ' '
```

The test pod's entrypoint should show `nono wrap` prepended to the original command.

## Configuration

The plugin TOML config is at `deploy/10-nono-nri.toml.example`. Key fields:
- `runtime_classes`: RuntimeClass names to intercept (default: `["nono-sandbox"]`)
- `nono_bin_path`: host path to nono binary (default: `/opt/nono-nri/nono`)
- `default_profile`: nono profile when annotation missing (default: `"default"`)

## Cleanup

```bash
kind delete cluster --name nono-test
```

## Customization

Override defaults with environment variables:
```bash
CLUSTER_NAME=my-cluster IMAGE=nono-nri:dev bash deploy/kind/deploy.sh
```
