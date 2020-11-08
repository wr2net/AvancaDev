FROM golang:1.15

WORKDIR /go/src/service_c
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN GOOS=linux go build

EXPOSE 9092

ENTRYPOINT [ "./service_c" ]