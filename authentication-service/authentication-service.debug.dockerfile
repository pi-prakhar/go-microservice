FROM golang:1.21.5-alpine

RUN apk add --no-cache go git

RUN go install github.com/go-delve/delve/cmd/dlv@v1.23.0

ENV PATH="/root/go/bin:${PATH}"

RUN which dlv

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN env GOOS=linux CGO_ENABLED=0 go build -gcflags="all=-N -l" -o authApp-debug ./cmd/api

EXPOSE 80 2345

CMD ["dlv", "exec", "./authApp-debug", "--headless=true", "--accept-multiclient" ,"--listen=:2345", "--log", "--api-version=2"]