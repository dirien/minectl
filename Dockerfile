# Dockerfile
FROM cgr.dev/chainguard/static@sha256:bda736ff75ec4b76ecd120d6ccf0b9428419cf4568a5d786a7d77ba23a9f6992
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
