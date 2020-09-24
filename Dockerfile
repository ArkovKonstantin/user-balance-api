FROM alpine
COPY ./bin/main /opt/main
COPY ./config /opt/config
RUN chmod +x /opt/main
WORKDIR /opt
RUN apk update && apk add ca-certificates