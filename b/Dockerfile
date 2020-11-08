FROM golang:1.15

WORKDIR /go/src/service_b
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN GOOS=linux go build

EXPOSE 9091

ENTRYPOINT [ "./service_b" ]