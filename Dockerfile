# Dockerfile
FROM alpine:3.14.3
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]