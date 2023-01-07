# Dockerfile
FROM cgr.dev/chainguard/static@sha256:8cfb0c91e14dc9a1859ba6c30eaeee0a6f5bec8d0a332fd6657760493d3165bd
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
