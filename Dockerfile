# Dockerfile
FROM cgr.dev/chainguard/static@sha256:853bfd4495abb4b65ede8fc9332513ca2626235589c2cef59b4fce5082d0836d
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
