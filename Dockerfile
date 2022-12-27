# Dockerfile
FROM cgr.dev/chainguard/static@sha256:4e108d39661f01865b821ed77e6a1792a6e910a5d8953dec08d258ad46fac044
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
