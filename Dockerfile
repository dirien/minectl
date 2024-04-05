# Dockerfile
FROM cgr.dev/chainguard/static@sha256:8665c8a9fcdab0f8afc09533ee23287c7870de26064d464a10e3baa52f337734
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
