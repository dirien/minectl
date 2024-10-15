# Dockerfile
FROM cgr.dev/chainguard/static@sha256:d07036a3beff43183f49bce5b2a0bd945f2ffe6e76f734ebd040059a40d371bc
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
