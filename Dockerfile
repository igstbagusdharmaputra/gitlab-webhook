FROM golang as builder
WORKDIR /build

COPY *.go go.mod go.sum  .
RUN go mod download  && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webhook *.go
FROM alpine
WORKDIR /app

COPY --from=builder /build/webhook .
CMD ["./webhook"]