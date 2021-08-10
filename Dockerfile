# Dockerfile
FROM alpine:3.14.1
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]