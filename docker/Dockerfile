FROM golang:1.14-alpine as build-img
LABEL maintainer jaimeenpatel22@gmail.com
ENV GO111MODULE=off

RUN apk update && apk add git && mkdir -p /home/tollPlaza
ARG TARGETSOURCEDIR=/go/src/automatedTollPlaze
COPY ${PWD}/ ${TARGETSOURCEDIR}/
WORKDIR ${TARGETSOURCEDIR}/
COPY ${PWD}/config/ /home/tollPlaza/config
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o /home/tollPlaza/tollPlaza cmd/*.go

FROM alpine:latest
WORKDIR /home
COPY --from=build-img /home/tollPlaza .
CMD ./tollPlaza
