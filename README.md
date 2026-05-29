# Zoo Bridge — fork template

Reference configuration + deployment scaffold for white-label bridge
consumers. **NOT** a separate deployment image. Zoo's own production
bridge runs the upstream `ghcr.io/luxfi/bridge` image directly with
tenant config supplied at runtime via env + ConfigMap.

## What lives here

| File | Purpose |
|---|---|
| `tenant.yaml` | Zoo Bridge config — brand, IAM endpoint `zoo.id`, KMS endpoint `kms.zoo.network`, MPC cluster, strict-pq profile, supported chains, basket allowlist, fee receiver, domain `bridge.zoo.network`, per-family release-pool sizing. |
| `contracts/tenant.json` + `contracts/Deploy.sh` | Wraps `@luxfi/standard v1.7.5+`'s `DeployTenant.s.sol` with the Zoo manifest. |
| `k8s/` | Deployment / Service / IngressRoute / ConfigMap manifests. |
| `tenant_test.go` | Validates `tenant.yaml` against the upstream `github.com/luxfi/bridge/pkg/tenant` schema. |
| `Dockerfile` | **Reference only.** Documents how a downstream fork that needs its own image (compliance bake-in like Liquidity) would assemble one. Zoo production does not build or publish this image. |

## Production deployment — config, not image

Zoo production pulls the canonical upstream image and injects tenant
config at deploy time via ConfigMap:

```yaml
spec:
  template:
    spec:
      containers:
        - name: bridge
          image: ghcr.io/luxfi/bridge:v1.1.40        # clean upstream semver
          args: ["--tenant-config", "/etc/bridge/tenant.yaml"]
          volumeMounts:
            - name: tenant-config
              mountPath: /etc/bridge
      volumes:
        - name: tenant-config
          configMap:
            name: zoo-bridge-tenant
```

ConfigMap rotation is hot-reloadable. Same image runs testnet, mainnet,
dev — environments differ only by ConfigMap.

## Why this repo exists

This is the **template** for any community white-label consumer:

```
git clone https://github.com/zooai/bridge
# edit tenant.yaml for your brand + endpoints
# point your k8s manifests at ghcr.io/luxfi/bridge:vX.Y.Z
# done.
```

If a downstream needs to bake config into a dedicated image (Liquidity
pattern: US ATS/BD/TA → region-locked GAR image, config baked in), the
`Dockerfile` here shows the minimal scaffold.

## Upstream pin

`luxfi/bridge` `v1.1.40` — image at `ghcr.io/luxfi/bridge:v1.1.40`.
Bump it in `k8s/deployment.yaml` to pick up new OSS features. Clean
semver — no tenant suffix.
