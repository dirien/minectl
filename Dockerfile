# Dockerfile
FROM ghcr.io/distroless/alpine-base@sha256:09ffedb199e159c97bd59338703e2fb82c73190be9b055558026c5ecdaf45b4a
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
