# Dockerfile
FROM alpine@sha256:d6d0a0eb4d40ef96f2310ead734848b9c819bb97c9d846385c4aca1767186cd4
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]