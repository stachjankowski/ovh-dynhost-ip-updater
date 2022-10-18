FROM golang:latest as builder

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /bin/ovh-dynhost-ip-updater

FROM alpine:latest
COPY --from=builder /bin/ovh-dynhost-ip-updater /bin/ovh-dynhost-ip-updater
ENTRYPOINT ["/bin/ovh-dynhost-ip-updater"]
