# Dockerfile
FROM ghcr.io/distroless/static@sha256:95361a439133b2148f0fbe07accb752c481ffa883ba8b3ee0a5c793e3ba747d3
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
