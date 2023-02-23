# Dockerfile
FROM cgr.dev/chainguard/static@sha256:0569f7d290a7892c62ebdd557f0c3c31ad3cf32ef49df8497cf65810add8a891
COPY minectl \
	/usr/bin/minectl
ENTRYPOINT ["/usr/bin/minectl"]
