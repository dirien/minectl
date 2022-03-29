# Dockerfile
FROM alpine@sha256:f22945d45ee2eb4dd463ed5a431d9f04fcd80ca768bb1acf898d91ce51f7bf04
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]