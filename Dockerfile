From golang:1.8.0-stretch AS builder

WORKDIR /go/src/github.com/pjgg/slackbot

COPY ./ ./
RUN make build
RUN ls /go/src/github.com/pjgg/slackbot

FROM frolvlad/alpine-glibc
RUN apk --no-cache upgrade && apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/pjgg/slackbot/bin/app-linux-amd64 ./app

RUN mkdir config
COPY config_example.yaml config/config.yaml

EXPOSE 8080

CMD /app
#ENTRYPOINT ["/app"]
