package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"github.com/wesleywillians/go-rabbitmq/queue"
	"io/ioutil"
	"log"
	"net/url"
	"github.com/hashicorp/go-retryablehttp"
	_ "github.com/wesleywillians/go-rabbitmq/queue"
	uuid "github.com/satori/go.uuid"
)

type Result struct {
	Status string
}

type Order struct {
	ID uuid.UUID
	Coupon string
	CcNumber string
}

func NewOrder() Order {
	return Order{
		ID: uuid.NewV4(),
	}
}

const (
	InvalidCouon = "invalid"
	ValidCoupon = "valid"
	ConnectionError = "connection error"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env", err)
	}
}

func main() {
	messageChannel := make(chan amqp.Delivery)

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	for msg := range messageChannel {
		process(msg)
	}
}

func process(msg amqp.Delivery) {

	order := NewOrder()
	json.Unmarshal(msg.Body, &order)

	resultCoupon := makeHttpCall("http://localhost:9092", order.Coupon)

	switch resultCoupon.Status {
	case InvalidCouon:
		log.Println("Order: ", order.ID, ": is invalid coupon!")

	case ConnectionError:
		msg.Reject(false)
		log.Println("Order: ", order.ID, ": could not process!")

	case ValidCoupon:
		log.Println("Order: ", order.ID, ": Processed")

	}
}

func makeHttpCall(urlMicroservice string, coupon string) Result {
	values := url.Values{}
	values.Add("coupon", coupon)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroservice, values)
	if err != nil {
		result :=  Result {Status: ConnectionError}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result
}