#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLUSTER_NAME="${CLUSTER_NAME:-nono-test}"
IMAGE="${IMAGE:-nono-nri:latest}"

echo "==> Creating Kind cluster '$CLUSTER_NAME'..."
kind create cluster --name "$CLUSTER_NAME" --config "$SCRIPT_DIR/cluster.yaml"

echo "==> Building nono-nri image..."
cd "$REPO_ROOT"
make docker-build IMAGE="$IMAGE"

echo "==> Loading image into Kind..."
kind load docker-image "$IMAGE" --name "$CLUSTER_NAME"

echo "==> Copying TOML config into Kind node..."
docker cp "$REPO_ROOT/deploy/10-nono-nri.toml.example" \
  "${CLUSTER_NAME}-control-plane:/etc/nri/conf.d/10-nono-nri.toml"

echo "==> Applying RuntimeClass..."
kubectl apply -f "$REPO_ROOT/deploy/runtimeclass.yaml"

echo "==> Applying DaemonSet..."
kubectl apply -f "$REPO_ROOT/deploy/daemonset.yaml"

echo "==> Waiting for DaemonSet rollout..."
kubectl rollout status daemonset/nono-nri -n kube-system --timeout=120s

echo ""
echo "==> Deployment complete!"
echo ""
echo "To test nono injection, run:"
echo "  kubectl apply -f $REPO_ROOT/deploy/test-pod.yaml"
echo "  kubectl wait --for=condition=ready pod/nono-test --timeout=60s"
echo "  kubectl exec nono-test -- cat /proc/1/cmdline | tr '\0' ' '"
echo ""
echo "To tear down:"
echo "  kind delete cluster --name $CLUSTER_NAME"
