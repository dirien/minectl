# Dockerfile
FROM cgr.dev/chainguard/static@sha256:110b6918893ea3df0eec04b2f469f3af07e5439900ed259076c55cefb1ec3965
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
