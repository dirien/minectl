# Dockerfile
FROM cgr.dev/chainguard/static@sha256:f0444e0c32aa0d22cca6de79e160306e9f24140fea0ac97c41789ee59814aa87
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
