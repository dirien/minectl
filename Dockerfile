# Dockerfile
FROM alpine@sha256:7580ece7963bfa863801466c0a488f11c86f85d9988051a9f9c68cb27f6b7872
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]