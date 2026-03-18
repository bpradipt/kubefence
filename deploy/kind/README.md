# Kind Local Test Clusters

Two cluster configurations are provided — one for each supported container runtime.

## Prerequisites

- Docker (running)
- [Kind](https://kind.sigs.k8s.io/) v0.20+
- kubectl
- Go 1.24+ (for building the plugin)
- `nono` binary at repo root (`./nono`) — dynamically linked against glibc + libdbus-1

## Supported Configurations

| Runtime | Image | SetArgs support | Notes |
|---------|-------|-----------------|-------|
| CRI-O 1.35+ | `quay.io/confidential-containers/kind-crio:v1.35.2` | ✓ | Full injection working |
| containerd 2.2.0+ | `kindest/node:v1.32.3` (ships containerd 2.0.3) | ✗ (2.0.x) | SetArgs fixed in containerd 2.2.0 |

> **Note on SetArgs:** `ContainerAdjustment.SetArgs()` is implemented in the NRI spec but was
> not applied in containerd ≤ 2.0.x or CRI-O ≤ 1.29.x due to a missing `AdjustArgs()` call in
> their vendored NRI runtime-tools library. CRI-O 1.35+ and containerd 2.2.0+ have the fix.

## Deploy

### CRI-O (recommended)

```bash
RUNTIME=crio bash deploy/kind/deploy.sh
```

This creates a cluster named `nono-crio` using `quay.io/confidential-containers/kind-crio:v1.35.2`.
A local Docker registry is started for image loading into CRI-O.

### containerd

```bash
RUNTIME=containerd bash deploy/kind/deploy.sh
```

This creates a cluster named `nono-containerd` using `kindest/node:v1.32.3`.
Image loading uses `ctr import` directly (workaround for `kind load docker-image` bug with containerd v2.x).

### Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `RUNTIME` | `containerd` | `containerd` or `crio` |
| `CLUSTER_NAME` | `nono-<runtime>` | Kind cluster name |
| `IMAGE` | `nono-nri:latest` | Plugin image tag |

## Run E2E Tests

After deploying with `deploy.sh`, run the 16-test e2e suite:

```bash
# CRI-O
RUNTIME=crio CLUSTER_NAME=nono-crio \
  REGISTRY_NAME=nono-nri-registry REGISTRY_PORT=5100 \
  bash deploy/kind/e2e.sh

# containerd
RUNTIME=containerd CLUSTER_NAME=nono-containerd \
  bash deploy/kind/e2e.sh
```

### E2E Test Coverage

| Test | What it verifies |
|------|-----------------|
| 1. Plugin connectivity | DaemonSet running, plugin registered with runtime |
| 2. Sandboxed pod injection | `process.args` modified, `/nono/nono` accessible, OCI bundle, state dir |
| 3. Non-sandboxed isolation | Non-sandboxed pods unaffected, no `/nono` mount |
| 4. State dir cleanup | State dir removed on pod deletion (`RemoveContainer`) |

### Expected Results by Runtime

| Test | CRI-O 1.35 | containerd 2.0.3 |
|------|-----------|-----------------|
| Plugin connectivity | ✓ | ✓ |
| process.args modified | ✓ | ✓ (nono exec'd via arg[0]) |
| /nono/nono accessible | ✓ | ✓ |
| OCI bundle process.args | ✓ | **✗** (SetArgs not applied, runtime limitation) |
| OCI bind mount | ✓ | ✓ |
| State dir metadata | ✓ | ✓ |
| Non-sandboxed isolation | ✓ | ✓ |
| State dir cleanup | ✓ | ✓ |

## Verify Manually

After deployment:

```bash
# Apply the test pod (uses nono-sandbox RuntimeClass)
kubectl apply -f deploy/test-pod.yaml

# Wait for it to be ready
kubectl wait --for=condition=ready pod/nono-test --timeout=60s

# Check /proc/1/cmdline — shows sleep (nono exec'd and replaced itself)
kubectl exec nono-test -- cat /proc/1/cmdline | tr '\0' ' '

# Check /nono/nono is bind-mounted
kubectl exec nono-test -- ls -la /nono/nono
```

## Cleanup

```bash
# CRI-O
kind delete cluster --name nono-crio
docker rm -f nono-nri-registry

# containerd
kind delete cluster --name nono-containerd
```
