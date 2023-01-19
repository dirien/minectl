# Dockerfile
FROM cgr.dev/chainguard/static@sha256:be7470ceeb0726d02f32c2ba2ae12cfe2025b2187eff36a53da30afd75fd2fa3
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
