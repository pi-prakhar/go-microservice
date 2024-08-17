FROM golang:1.21.5-alpine

RUN mkdir /app

COPY authApp /app

WORKDIR /app

CMD ["./authApp"]