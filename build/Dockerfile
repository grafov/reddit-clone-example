### STAGE 1: Build  ----------------------------------------------------------------------------------------------------
FROM golang:1.15 as builder

WORKDIR /go/src/redditclone

COPY . .
RUN make build

### STAGE 2: Setup  ----------------------------------------------------------------------------------------------------
FROM scratch

VOLUME /tmp

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/redditclone/web/* /web/
COPY --from=builder /go/src/redditclone/storage/migrations/* /storage/migrations/
COPY --from=builder /go/src/redditclone/rcapp /rcapp

ENTRYPOINT ["/rcapp"]
