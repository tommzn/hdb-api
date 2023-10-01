FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.19 as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG Version
ARG GitCommit

ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /go/build

# Cache the download before continuing
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY .  .

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -ldflags "-s -w -X github.com/inlets/inlets-operator/pkg/version.Release=${Version} -X github.com/inlets/inlets-operator/pkg/version.SHA=${GitCommit}" \
  -a -installsuffix cgo -o build_artifact_bin .

FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.source=https://github.com/inlets/inlets-operator

WORKDIR /go

COPY --from=builder /go/build/build_artifact_bin hdb-bin
USER nonroot:nonroot

ENTRYPOINT ["/go/hdb-bin"]
