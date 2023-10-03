FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.19-alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG Version
ARG GitCommit

ENV CGO_ENABLED=1
ENV GO111MODULE=on

RUN apk add build-base librdkafka-dev pkgconf

WORKDIR /go/build

# Cache the download before continuing
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download


COPY .  .

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go test -v ./...

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
  go build -tags dynamic -v -o build_artifact_bin

FROM --platform=${BUILDPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.source=https://github.com/tommzn/hdb-api

WORKDIR /go

COPY --from=builder /go/build/build_artifact_bin hdb-bin
USER nonroot:nonroot

EXPOSE 8080
ENTRYPOINT ["/go/hdb-bin"]
