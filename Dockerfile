# Dockerfile
FROM cgr.dev/chainguard/static@sha256:368a37f7a803dba397b2f1a199dae66ee94746ffdd0b5c61626d48e261bef321
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
