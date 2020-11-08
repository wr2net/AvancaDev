FROM golang:1.15

WORKDIR /go/src/service_a
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN GOOS=linux go build

EXPOSE 9090

ENTRYPOINT [ "./service_a" ]