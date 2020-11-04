FROM golang:alpine3.10
WORKDIR /root
COPY ./ ./
RUN CGO_ENABLED=0 go build -o plugin/hdwallet *.go
CMD ["sleep","10m"]