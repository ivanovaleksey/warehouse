FROM golang:1.22-alpine AS builder

RUN apk add --update make

WORKDIR /build

COPY go.* .
RUN go mod download

COPY Makefile Makefile
RUN make install-proto-tools

COPY api/ api/
COPY cmd/ cmd/
COPY config/ config/
COPY db/ db/
COPY internal/ internal/
RUN make build-server
RUN make build-seed

FROM alpine

WORKDIR /app

COPY --from=builder /build/bin/server .
COPY --from=builder /build/bin/seed .
COPY config/ config/
COPY db/ db/
COPY cmd/seed/data seeddata/

CMD ["/app/server", "-config", "config/config.docker.yaml"]
