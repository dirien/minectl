# Dockerfile
FROM cgr.dev/chainguard/static@sha256:0bac1022a66a2989609e8ff7cdacbc35d1a8692670bd05ae305b90e1c0eba337
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
