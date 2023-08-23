FROM alpine:latest

WORKDIR /app

ADD build.tar /app
ADD webssh_target /app/build

CMD ["/app/build/webssh_target"]