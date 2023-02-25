# Dockerfile
FROM cgr.dev/chainguard/static@sha256:ce53292a08ad6a82f89183887b30fbcc7dea84bf3a3155e70516f323d699e83a
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
