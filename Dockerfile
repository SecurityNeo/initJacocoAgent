FROM alpine:3.13.7

ADD initJacocoAgent /initJacocoAgent
ENTRYPOINT ["./initJacocoAgent"]