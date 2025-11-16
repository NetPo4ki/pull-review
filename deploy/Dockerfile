FROM golang:1.24 AS builder
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app /app/app
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app/app"]