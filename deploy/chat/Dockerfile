FROM golang:1.19-alpine as builder

RUN mkdir -p "/app"
WORKDIR /app
COPY .. .

RUN go build ./cmd/web

FROM alpine:3.16.2

RUN mkdir -p "/tmp"
WORKDIR "/tmp"

ARG COMMAND

COPY --from=builder /app/web .
RUN chmod +x ./web

ENTRYPOINT ./web








