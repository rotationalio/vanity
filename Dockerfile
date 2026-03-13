# Dynamic Builds
ARG BUILDER_IMAGE=dhi.io/golang:1.25-debian13-dev
ARG FINAL_IMAGE=dhi.io/golang:1.25-debian13

# Build stage
FROM --platform=${BUILDPLATFORM} ${BUILDER_IMAGE} AS builder

# Build Args
ARG GIT_REVISION=""
ARG BUILD_DATE=""

# Platform args
ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM

# Ensure ca-certificates are up to date
RUN update-ca-certificates

# Use modules for dependencies
WORKDIR $GOPATH/src/go.rtnl.ai/vanity

COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# Copy package
COPY . .

# Build binary
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v \
    -ldflags="-X 'go.rtnl.ai/vanity.GitVersion=${GIT_REVISION}' -X 'go.rtnl.ai/vanity.BuildDate=${BUILD_DATE}'" \
    -o /go/bin/vanityd \
    ./cmd/vanityd

# Final Stage
FROM --platform=${BUILDPLATFORM} ${FINAL_IMAGE} AS final

LABEL maintainer="Rotational Labs <support@rotational.io>"
LABEL description="Rotational Vanity URLs server for go modules"

# Copy the binary to the production image from the builder stage
COPY --from=builder /go/bin/vanityd /usr/local/bin/vanityd

CMD [ "/usr/local/bin/vanityd", "serve" ]
