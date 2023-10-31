# Dockerfile
FROM cgr.dev/chainguard/static@sha256:d3465871ccaba3d4aefe51d6bb2222195850f6734cbbb6ef0dd7a3da49826159
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
