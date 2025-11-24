# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Compile statically to keep the runtime image small
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/bin/nexivent ./internal

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app

ENV PORT=8080

COPY --from=builder /app/bin/nexivent /app/nexivent

EXPOSE 8080

CMD ["/app/nexivent"]
