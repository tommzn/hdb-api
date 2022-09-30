FROM alpine:latest

WORKDIR /go

COPY --chmod=0755 build_artifact_bin hdb-bin

ENTRYPOINT ["/go/hdb-bin"]
