FROM golang:latest

WORKDIR /usr/pollie
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/main.go

FROM alpine:latest as runner
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /usr/pollie/app .
COPY --from=0 /usr/pollie/config_docker.yml ./config.yml
EXPOSE 6000
CMD ["./app"]  