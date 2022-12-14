# Dockerfile
FROM cgr.dev/chainguard/static@sha256:34e0f01926aa86f932fdd6e5d8f4e24e186a3e55b46966269b8dc78dfcd7353a
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
