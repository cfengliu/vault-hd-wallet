FROM golang:1.14-alpine
RUN apk add --update alpine-sdk
RUN apk update && apk add git openssh gcc musl-dev linux-headers

WORKDIR /root
COPY ./ ./

RUN CGO_ENABLED=1 GOOS=linux go build -a -v -i -o plugin/hdwallet *.go
CMD ["sleep","10m"]