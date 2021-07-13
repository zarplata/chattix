FROM golang:latest as builder

ARG MESSENGER=slack
ENV MESSENGER ${MESSENGER}

WORKDIR /go/src/github.com/zarplata/chattix
RUN go get -u github.com/golang/dep/...

COPY . .
RUN CGO_ENABLED=0 make service

FROM alpine:latest
RUN apk add --update bash ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /bin/
COPY --from=builder /go/src/github.com/zarplata/chattix/.out/* chattixd

CMD ["/bin/chattixd"]
