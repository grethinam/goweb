FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./app /go/src/github.com/user/web/app
WORKDIR /go/src/github.com/user/web/app

RUN go get ./
RUN go build

EXPOSE 8080