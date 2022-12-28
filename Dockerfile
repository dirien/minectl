# Dockerfile
FROM cgr.dev/chainguard/static@sha256:7f89515b06306eab84e8d80e154779ba2e047dd17a5ceb73f4c0f84e8c210d20
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
