FROM docker.io/golang:1.21.6-alpine AS build
LABEL MAINTAINER github.com/arizon-dread

WORKDIR /usr/local/go/src/github.com/arizon-dread/webdig-backend
COPY businessLayer ./businessLayer
COPY models ./models
COPY config ./config
COPY api.go main.go go.mod go.sum ./

RUN apk update && apk add --no-cache git
RUN go build -v -o /usr/local/bin/webdig-backend/ ./...


FROM docker.io/golang:1.21.6-alpine AS final
WORKDIR /go/bin
ENV GIN_MODE=release
#RUN apk add --no-cache libc6-compat musl-dev
COPY --from=build /usr/local/bin/webdig-backend/ /go/bin/
EXPOSE 8080

ENTRYPOINT [ "./webdig-backend" ]
