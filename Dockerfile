FROM golang:1.14.2-alpine3.11

RUN apk update && apk --no-cache add ca-certificates

RUN apk update && apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Kolkata /etc/localtime
RUN echo "Asia/Kolkata" > /etc/timezone
ARG BUILD_FOR=prod

RUN apk update && apk add git
RUN apk update && apk add curl
RUN mkdir -p /covid-app/
ADD . /covid-app/

ENV SERVICE=covid-app
WORKDIR /covid-app
COPY /config/$BUILD_FOR/config.go ./config/config.go
WORKDIR /covid-app

RUN go build -o server .

EXPOSE 8080

# CMD gunicorn --bind 0.0.0.0:$PORT wsgi
CMD ["/covid-app/server"]