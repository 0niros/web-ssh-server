FROM ubuntu:latest

ADD build.tar /app
ADD webssh_target /app/build
WORKDIR /app/build

CMD ["/app/build/webssh_target"]
