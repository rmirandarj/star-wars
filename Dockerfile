FROM golang:1.16.5-alpine3.14 as BUILDER
RUN apk update && apk add --no-cache git ca-certificates tzdata openssh make
RUN adduser -D -g '' appuser
WORKDIR /app
COPY . .
RUN make -S build

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
WORKDIR /app
COPY --from=builder /app/build/server ./server
USER appuser
ENTRYPOINT ["/app/server"]