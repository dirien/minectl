# Dockerfile
FROM cgr.dev/chainguard/static@sha256:67a1b00e0134e2b3a614c7198a26f7deed9d11b7acad4d52c79c0cfd47a2eae7
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
