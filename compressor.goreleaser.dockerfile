FROM gcr.io/distroless/static-debian11
COPY compressor /compressor
ENTRYPOINT ["/compressor"]
