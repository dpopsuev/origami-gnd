# GND (Gather & Diffuse) — Papercup circuit server.
# Exposes tree → search → read → synthesize via start_circuit protocol.
#
# Build context must include sibling repos:
#   docker build -t origami-gnd -f rh-dsr/Dockerfile .   (from workspace root)
#
# Run: docker run -p 9100:9100 origami-gnd --port 9100 --drivers git,docs

FROM golang:1.24 AS builder
WORKDIR /src
COPY rh-dsr/go.mod rh-dsr/go.sum ./rh-dsr/
COPY origami/go.mod origami/go.sum ./origami/
RUN cd rh-dsr && \
    go mod edit \
        -replace github.com/dpopsuev/origami=../origami && \
    go mod download
COPY rh-dsr/ ./rh-dsr/
COPY origami/ ./origami/
RUN cd rh-dsr && CGO_ENABLED=0 go build -o /gnd-serve ./cmd/serve

FROM gcr.io/distroless/static-debian12
COPY --from=builder /gnd-serve /gnd-serve
ENTRYPOINT ["/gnd-serve"]
EXPOSE 9100
