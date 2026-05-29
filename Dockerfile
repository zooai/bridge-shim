# Zoo Bridge production image.
FROM ghcr.io/luxfi/bridge:v1.1.40

COPY tenant.yaml /etc/bridge/tenant.yaml

ENTRYPOINT ["/usr/local/bin/bridge", "--tenant-config", "/etc/bridge/tenant.yaml"]
