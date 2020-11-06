package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"github.com/hashicorp/go-retryablehttp"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}

func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

type Result struct {
	Status string
}

var coupons Coupons

func main() {
	coupon := Coupon{
		Code: "abc",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	valid := coupons.Check(coupon)

	result := Result{Status: valid}

	resultServe := makeHttpCall("http://localhost:9093/%22")
	jsonResultserve, err := json.Marshal(resultServe)
	if err != nil {
		log.Fatal("Failed in converting json")
	}
	log.Println(string(jsonResultserve))

	fmt.Fprintf(w, string(jsonResultserve))

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Failed in converting json")
	}

	fmt.Fprintf(w, string(jsonResult))
}

func makeHttpCall(uriService string) Result {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	values := url.Values{}

	res, err := retryClient.PostForm(uriService, values)
	if err != nil {
		result := Result{Status: "Service offline"}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Fail in processing result")
	}

	result := Result{}
	json.Unmarshal(data, &result)
	return result
}