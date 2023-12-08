# Dockerfile
FROM cgr.dev/chainguard/static@sha256:6ddc52909b55f716c3154a28efab0c7ffb7ef741aeb390a7a431e9f3001ec010
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
