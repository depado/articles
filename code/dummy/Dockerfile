# Build step
FROM golang:latest AS build

RUN mkdir -p $GOPATH/src/github.com/Depado/dummy
ADD . $GOPATH/src/github.com/Depado/dummy
WORKDIR $GOPATH/src/github.com/Depado/dummy
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure -vendor-only
RUN CGO_ENABLED=0 go build -o /dummy

# Final step
FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
RUN apk add --update tzdata
RUN rm -rf /var/cache/apk/*

COPY --from=build /dummy /home/
ENV TZ=Europe/Paris
WORKDIR /home
ENTRYPOINT ./dummy
EXPOSE 8080