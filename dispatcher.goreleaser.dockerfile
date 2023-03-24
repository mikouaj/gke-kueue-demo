FROM gcr.io/distroless/static-debian11
COPY dispatcher /dispatcher
ENTRYPOINT ["/dispatcher"]
