# Dockerfile
FROM cgr.dev/chainguard/static@sha256:bee65530377f28c6ef1090a7b1878c08ba97181c94b30376e9ac37bb8ad6ddb4
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
