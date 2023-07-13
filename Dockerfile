FROM golang:1.20 as builder

WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download

COPY src/** /app
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o app .

FROM gcr.io/distroless/static

COPY --from=builder /app/app /app

ENTRYPOINT ["/app"]
