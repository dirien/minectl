# Dockerfile
FROM gcr.io/distroless/base-debian10
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]