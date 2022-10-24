FROM golang:1.19.2-alpine3.16 as builder

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /bin/ovh-dynhost-ip-updater


FROM alpine:3.16.2

COPY --from=builder /bin/ovh-dynhost-ip-updater /bin/ovh-dynhost-ip-updater

USER 1001

ENTRYPOINT ["/bin/ovh-dynhost-ip-updater"]
