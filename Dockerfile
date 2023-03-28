# Dockerfile
FROM cgr.dev/chainguard/static@sha256:e3bfa4cfcee80b1a32dcf09b8c33d0255ce502a3c877d19f77af09d591373d78
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
