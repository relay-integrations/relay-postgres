FROM golang:1.14-alpine AS builder
ENV CGO_ENABLED 0
WORKDIR /build
COPY . .
RUN go build -a -o "/usr/local/bin/relay-postgres-query" "./cmd/postgres-query"

FROM alpine:3
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=builder "/usr/local/bin/relay-postgres-query" "/usr/local/bin/reqlay-postgres-query"
ENTRYPOINT []
CMD ["/usr/local/bin/relay-postgres-query"]

LABEL "org.opencontainers.image.title"="Issue a query to Postgres."
LABEL "org.opencontainers.image.description"="This step issues a query to Postgres and sets the results as an output."
LABEL "com.puppet.nebula.sdk.version"="v1"
