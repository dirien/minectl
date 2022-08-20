# Dockerfile
FROM ghcr.io/distroless/alpine-base@sha256:6b798ef68d029c2311e461b1dfdedfbb8cff4784a31088811960d23745336caf
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
