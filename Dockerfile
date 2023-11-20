# Dockerfile
FROM cgr.dev/chainguard/static@sha256:676e989769aa9a5254fbfe14abb698804674b91c4d574bb33368d87930c5c472
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
