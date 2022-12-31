# Dockerfile
FROM cgr.dev/chainguard/static@sha256:bddbb08d380457157e9b12b8d0985c45ac1072a1f901783a4b7c852e494967d8
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
