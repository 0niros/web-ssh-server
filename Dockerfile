FROM alpine:latest

WORKDIR /app

ADD build.tar /app
ADD webssh_target /app

CMD ["/app/webssh_target"]