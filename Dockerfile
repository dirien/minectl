# Dockerfile
FROM alpine:3.14.2
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]