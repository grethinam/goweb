FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./app /go/src/goweb/app
WORKDIR /go/src/goweb/app

RUN go get ./
RUN go build

ENTRYPOINT /go/bin/app

EXPOSE 8080