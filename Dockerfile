# Dockerfile
FROM cgr.dev/chainguard/static@sha256:3824d30a74e13d80af2ccb852516b6b7b2684b75d9f100b79d13863ae35626b5
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
