FROM golang:1.22-alpine

RUN apk add --update make protobuf

WORKDIR /build

COPY go.* .
RUN go mod download

COPY Makefile Makefile
RUN make install-tools

COPY api/ api/
COPY cmd/ cmd/
COPY config/ config/
COPY tests/config.yaml config/config.test.yaml
COPY db/ db/
COPY internal/ internal/
RUN make generate
