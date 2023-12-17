# Dockerfile
FROM cgr.dev/chainguard/static@sha256:ac42c519db36f29a71f26f7ec915bec198a84ff76b341eb8fcd9c4e73e31d356
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
