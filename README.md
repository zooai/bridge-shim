# Zoo Bridge

Zoo's tenant configuration bundle for the canonical Lux bridge
(`bridge.zoo.network`). Runs the standard `ghcr.io/luxfi/bridge`
image; everything Zoo-specific is config, not code.

> **Pattern: declarative Bridge SDK.** Same idea as the `@luxfi/bridge`
> browser SDK (`mountBridge(config)`) and the upstream `@luxfi/exchange`
> SDK — one canonical implementation, every tenant declares its
> identity via config. No fork, no shim, no per-tenant image.
> **Rename target: `zooai/bridge-shim` → `zooai/bridge`.**

## What lives here

| File | Purpose |
|---|---|
| `tenant.yaml` | Zoo tenant config (brand, IAM endpoint `zoo.id`, KMS endpoint `kms.zoo.network`, MPC cluster, strict-pq profile, supported chains, basket allowlist, fee receiver, domain `bridge.zoo.network`, per-family release-pool sizing). |
| `contracts/tenant.json` + `contracts/Deploy.sh` | Wraps `@luxfi/standard v1.7.5+`'s `DeployTenant.s.sol` with the Zoo manifest — drops the Zoo BridgeV4 + cosigner contracts on-chain. |
| `k8s/` | Deployment / Service / IngressRoute / ConfigMap manifests targeting `bridge.zoo.network`, pulling `ghcr.io/luxfi/bridge:vX.Y.Z` directly. |
| `tenant_test.go` | Validates `tenant.yaml` against `github.com/luxfi/bridge/pkg/tenant.Config`. CI runs on every PR — same schema gate every Lux tenant uses. |

## Deploy

```bash
# 1. Validate config locally.
go test ./...

# 2. Apply k8s manifests — pulls upstream image, mounts tenant.yaml.
kubectl apply -k k8s/

# 3. Deploy contracts.
ZOO_PRIVATE_KEY=0x... ./contracts/Deploy.sh https://rpc.zoo.network
```

The `kubectl apply` step deploys `ghcr.io/luxfi/bridge:v1.1.40`
unchanged. Zoo identity comes from the ConfigMap created from
`tenant.yaml` — not from a Zoo-branded image. Bumping upstream is one
`vX.Y.Z` change in `k8s/deployment.yaml`.

## Why declarative SDK pattern

```
                    ┌─────────────────────────────┐
                    │ ghcr.io/luxfi/bridge:vX.Y.Z │ ← canonical OSS image
                    │   (brand-neutral binary)    │
                    └────────────┬────────────────┘
                                 │ runtime config
        ┌────────────────────────┼────────────────────────┐
        ▼                        ▼                        ▼
   zoo-bridge-tenant      hanzo-bridge-config    (no separate repo
   (this repo's            (~/work/hanzo/platform   needed — config
   tenant.yaml →           ConfigMap → applied      lives in platform
   ConfigMap)              direct to k8s)           infra)
        ▼                        ▼                        ▼
   bridge.zoo.network     bridge.hanzo.network    bridge.{anyone}.tld
```

Exception: regulated downstream variants (US ATS/BD/TA) may bake
config into a region-locked GAR image because compliance requires
no cross-region config mutations. Everywhere else, runtime config wins.

## Image tagging convention

Clean upstream semver. The image tag IS `ghcr.io/luxfi/bridge:v1.1.40`
— there is no Zoo-tagged image. Bump it in `k8s/deployment.yaml` to
pick up new OSS features.

(Historical: the `zooai/bridge` GHCR package shipped 12 versions
between 2026-03-31 and 2026-05-06 under an earlier per-tenant-image
convention. That package was deleted in the 2026-05-29 cleanup.
Runtime config is the path going forward.)

## Fork starter

If you want to white-label the Lux bridge for your own deployment,
clone this repo, edit `tenant.yaml` for your brand + endpoints, point
your `k8s/deployment.yaml` at `ghcr.io/luxfi/bridge:vX.Y.Z`, deploy.
That's the entire white-label flow.
