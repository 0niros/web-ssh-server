FROM alpine:latest

WORKDIR /app

ADD build.tar /app/build
ADD webssh_target /app

CMD ["/app/webssh_target"]