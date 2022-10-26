FROM alpine:3.13.7

ADD initSkywalkingAgent /initSkywalkingAgent
ENTRYPOINT ["./initSkywalkingAgent"]