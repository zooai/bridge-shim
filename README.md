# Zoo Bridge Shim

White-label tenant shim for the canonical bridge — `bridge.zoo.network`.

The runtime binary is `ghcr.io/luxfi/bridge` (the OSS upstream). This
repo ships:

- `tenant.yaml` — the declarative Zoo Bridge configuration (brand,
  IAM endpoint `zoo.id`, KMS endpoint `kms.zoo.network`, MPC cluster,
  strict-pq profile, supported chains, basket allowlist, fee receiver,
  domain `bridge.zoo.network`, per-family release-pool sizing).
- `Dockerfile` — overlays `tenant.yaml` on top of the upstream image.
- `contracts/tenant.json` + `contracts/Deploy.sh` — wraps
  `@luxfi/standard v1.7.5+`'s `DeployTenant.s.sol`.
- `k8s/` — Deployment, Service, IngressRoute, ConfigMap manifests.
- `tenant_test.go` — validates `tenant.yaml` against the upstream
  schema.

## Deploy

```bash
go test ./...
docker build -t ghcr.io/zooai/bridge-shim:v0.1.0 .
kubectl apply -k k8s/
ZOO_PRIVATE_KEY=0x... ./contracts/Deploy.sh https://rpc.zoo.network
```

## Upstream pins

- `ghcr.io/luxfi/bridge`: v2.0.0
- `@luxfi/standard`: v1.7.5
