FROM gcr.io/distroless/static
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/hookfeed /
ENTRYPOINT ["/hookfeed", "serve"]
