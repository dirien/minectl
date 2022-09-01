# Dockerfile
FROM ghcr.io/distroless/static@sha256:50c3d52ecfaed0112b6743df9dde26aaaf61a4d0b89ca167ade636678c793eb2
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
