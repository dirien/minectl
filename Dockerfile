# Dockerfile
FROM cgr.dev/chainguard/static@sha256:ddfcc031a62d7b6b97cbe120b7dbc051f9b64fce140acb33db29fb37a2d3889e
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
