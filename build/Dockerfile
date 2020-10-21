## STAGE 1: Build  ----------------------------------------------------------------------------------------------------
FROM golang:1.15-alpine as builder

### Add certs for possible HTTP clients
RUN apk add git make
RUN apk add --update --no-cache ca-certificates

WORKDIR /go/src/redditclone

### Cache modules
COPY go.mod go.sum ./
RUN go mod download

### Build app
COPY . .
RUN make build

## STAGE 2: Setup  ----------------------------------------------------------------------------------------------------
FROM scratch

VOLUME /tmp

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/redditclone/web/* /web/
COPY --from=builder /go/src/redditclone/storage/migrations/* /storage/migrations/
COPY --from=builder /go/src/redditclone/rcapp /rcapp

ENTRYPOINT ["/rcapp"]
