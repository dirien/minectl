# Dockerfile
FROM alpine:3.15.0
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]