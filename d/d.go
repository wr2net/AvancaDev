package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	result := Result{Status: "valid"}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Fail in convertint json")
	}

	fmt.Fprint(w, string(jsonResult))
}