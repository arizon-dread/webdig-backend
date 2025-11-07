FROM docker.io/golang:1.25-alpine AS build
LABEL MAINTAINER=github.com/arizon-dread

WORKDIR /usr/local/go/src/github.com/arizon-dread/webdig-backend
COPY . .

RUN apk update && apk add --no-cache git
RUN go build -v -o /usr/local/bin/webdig-backend/ ./...


FROM docker.io/golang:1.25-alpine AS final
WORKDIR /go/bin
ENV GIN_MODE=release
#RUN apk add --no-cache libc6-compat musl-dev
COPY --from=build /usr/local/bin/webdig-backend/ /go/bin/
EXPOSE 8080

ENTRYPOINT [ "./webdig-api" ]
