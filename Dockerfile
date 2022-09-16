FROM alpine:latest

RUN mkdir /app

COPY httpecho /app

ENTRYPOINT [ "/app/httpecho" ]