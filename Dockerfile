# Dockerfile
FROM cgr.dev/chainguard/static@sha256:eae79cd3ef653a6233fa1f87d3643cb3ff279eebf78fd8ff83d1a44255a06628
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
