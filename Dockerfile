# Dockerfile
FROM alpine@sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]