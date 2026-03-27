FROM docker.io/golang:1.26-alpine AS build
LABEL MAINTAINER=github.com/arizon-dread

WORKDIR /usr/local/go/src/github.com/arizon-dread/webdig-backend
COPY . .
RUN rm -r ./cmd/webdig-cli

RUN apk update && apk add --no-cache git
RUN go build -v -o /usr/local/bin/webdig-backend/ ./...


FROM dhi.io/alpine-base:3.23 AS final
WORKDIR /go/bin
ARG VERSION
ENV GENERAL_VERSION=${VERSION}
#RUN apk add --no-cache libc6-compat musl-dev
COPY --from=build /usr/local/bin/webdig-backend/ /go/bin/
EXPOSE 8080
ENTRYPOINT [ "./webdig-api" ]
