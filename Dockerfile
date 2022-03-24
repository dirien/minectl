# Dockerfile
FROM alpine@sha256:ceeae2849a425ef1a7e591d8288f1a58cdf1f4e8d9da7510e29ea829e61cf512
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]