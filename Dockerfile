# Dockerfile
FROM cgr.dev/chainguard/static@sha256:89c6f614cab203d9f77d8636231f40cd9cf89cac2ca5d353cd75186a1ce4c5f9
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
