# Kind Local Test Clusters

Two cluster configurations are provided — one for each supported container runtime.

## Prerequisites

- Docker (running)
- [Kind](https://kind.sigs.k8s.io/) v0.20+
- kubectl
- Go 1.24+ (for building the plugin)
- `nono` binary at repo root (`./nono`) — dynamically linked against glibc + libdbus-1

## Supported Configurations

| Runtime | Kind node image | containerd/CRI-O version | SetArgs support |
|---------|----------------|--------------------------|-----------------|
| containerd | `kindest/node:v1.35.1` | containerd 2.2.1 | ✓ |
| CRI-O | `quay.io/confidential-containers/kind-crio:v1.35.2` | CRI-O 1.35 | ✓ |

> **Note on SetArgs:** `ContainerAdjustment.SetArgs()` requires containerd ≥ 2.2.0 or
> CRI-O ≥ 1.35. Earlier versions had a missing `AdjustArgs()` call in their vendored
> NRI runtime-tools library and silently ignored args modifications.

## Deploy

### Using Make (recommended)

```bash
# containerd (default)
make kind-e2e

# CRI-O
make kind-e2e RUNTIME=crio

# Deploy only (keep cluster running for manual inspection)
make kind-up
make kind-up RUNTIME=crio

# Run tests against an existing cluster
make kind-test

# Tear down
make kind-down
make kind-down RUNTIME=crio
```

### Using the script directly

```bash
# containerd
RUNTIME=containerd bash deploy/kind/deploy.sh

# CRI-O
RUNTIME=crio bash deploy/kind/deploy.sh
```

### Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `RUNTIME` | `containerd` | `containerd` or `crio` |
| `CLUSTER_NAME` | `nono-<runtime>` | Kind cluster name |
| `IMAGE` | `nono-nri:latest` | Plugin image tag |
| `KATA` | `false` | Install Kata Containers (`true`/`false`) |
| `REGISTRY_NAME` | `nono-nri-registry` | Local registry container name (crio only) |
| `REGISTRY_PORT` | `5100` | Local registry port on the host (crio only) |

## Run E2E Tests

After deploying with `deploy.sh` or `make kind-up`:

```bash
# containerd
make kind-test
# or
RUNTIME=containerd CLUSTER_NAME=nono-containerd bash deploy/kind/e2e.sh

# CRI-O
make kind-test RUNTIME=crio
# or
RUNTIME=crio CLUSTER_NAME=nono-crio \
  REGISTRY_NAME=nono-nri-registry REGISTRY_PORT=5100 \
  bash deploy/kind/e2e.sh
```

### E2E Test Coverage

| Test | What it verifies |
|------|-----------------|
| 1. Plugin connectivity | DaemonSet running, plugin registered with runtime |
| 2. Sandboxed pod injection | `process.args` modified, `/nono/nono` accessible, OCI bundle args + mount, state dir written |
| 3. Non-sandboxed isolation | Non-sandboxed pods unaffected, no `/nono` mount |
| 4. State dir cleanup | State dir removed on pod deletion (`RemoveContainer`) |
| 5. Kata + nono | nono injection inside a QEMU/KVM micro-VM (skipped when `KATA=false`) |

### Expected Results

| Test | containerd 2.2.1 | CRI-O 1.35 |
|------|-----------------|------------|
| Plugin connectivity | ✓ | ✓ |
| process.args modified | ✓ | ✓ |
| /nono/nono accessible | ✓ | ✓ |
| OCI bundle process.args | ✓ | ✓ |
| OCI bind mount | ✓ | ✓ |
| State dir metadata | ✓ | ✓ |
| Non-sandboxed isolation | ✓ | ✓ |
| State dir cleanup | ✓ | ✓ |
| Kata + nono | skipped (KATA=false) | skipped (KATA=false) |

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
# containerd
make kind-down
# or: kind delete cluster --name nono-containerd

# CRI-O
make kind-down RUNTIME=crio
# or:
kind delete cluster --name nono-crio
docker rm -f nono-nri-registry
```
