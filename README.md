# Zoo Bridge

Tenant deployment of the canonical bridge — `bridge.zoo.network`.

> Not a "shim" — Zoo Bridge tenant deployment. Runtime binary IS
> `ghcr.io/luxfi/bridge` (the OSS upstream), composed at deploy time
> via tenant.yaml. **Rename target: `zooai/bridge-shim` → `zooai/bridge`.**

This repo ships:

- `tenant.yaml` — declarative Zoo Bridge configuration (brand, IAM
  endpoint `zoo.id`, KMS endpoint `kms.zoo.network`, MPC cluster,
  strict-pq profile, supported chains, basket allowlist, fee receiver,
  domain `bridge.zoo.network`, per-family release-pool sizing).
- `Dockerfile` — overlays `tenant.yaml` on top of the upstream image.
- `contracts/tenant.json` + `contracts/Deploy.sh` — wraps
  `@luxfi/standard v1.7.5+`'s `DeployTenant.s.sol`.
- `k8s/` — Deployment, Service, IngressRoute, ConfigMap manifests.
- `tenant_test.go` — validates `tenant.yaml` against the upstream schema.

## Image tagging convention

Tag the tenant image by the upstream version it composes:

```
ghcr.io/zooai/bridge:v1.1.40-zoo
```

NOT independent semver. The Dockerfile `FROM` pin IS the contract.
Legacy `v0.1.x` / `v0.2.x` tags remain in history; new deploys use
the `vX.Y.Z-zoo` form.

## Deploy

```bash
go test ./...
docker build -t ghcr.io/zooai/bridge:v1.1.40-zoo .
kubectl apply -k k8s/
ZOO_PRIVATE_KEY=0x... ./contracts/Deploy.sh https://rpc.zoo.network
```

## Upstream pins

- `ghcr.io/luxfi/bridge`: v1.1.40
- `@luxfi/standard`: v1.7.5
