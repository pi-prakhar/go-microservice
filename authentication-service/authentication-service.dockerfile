FROM alpine:latest

RUN apk add --no-cache go git

RUN go install github.com/go-delve/delve/cmd/dlv@latest

ENV PATH="/root/go/bin:${PATH}"

RUN mkdir /app

COPY authApp /app

EXPOSE 8081 2345

RUN which dlv

WORKDIR /app

CMD ["dlv", "exec", "./authApp", "--headless", "--listen=:2345", "--log", "--api-version=2"]