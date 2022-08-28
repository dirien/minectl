# Dockerfile
FROM ghcr.io/distroless/static@sha256:a438a0a33b67a4d6348457d8b57e7e547a894b5d1d0dae133666e27d92a1fa14
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
