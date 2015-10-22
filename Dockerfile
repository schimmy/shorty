# shorty URL shortener service
FROM golang:1.4

ENV service "shorty"
ENV dir "/go/src/github.com/Clever/$service"

RUN mkdir -p "$dir"
ADD . "$dir"
WORKDIR "$dir"

RUN go get ./...
RUN go build

CMD ["./shorty"]
