FROM alpine:latest

RUN mkdir /app

COPY frontApp /app

EXPOSE 8081

CMD [ "/app/frontApp"]