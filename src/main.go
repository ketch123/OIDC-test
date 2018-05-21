package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

const (
	rs2Letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	oauthURL      = "https://accounts.google.com/o/oauth2/v2/auth"
	response_type = "code"
	scope         = "email"
	redirect_uri  = "http://localhost/auth"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", oauth)
	router.Handle("/", router)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", "8080"), router); err != nil {
		log.Fatal("err: %v", err)
	}
}

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rs2Letters[rand.Intn(len(rs2Letters))]
	}
	return string(b)
}

func oauth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("oauth started...")
	state := RandString(15)
	nonce := RandString(20)

	values := url.Values{}
	values.Add("client_id", os.Getenv("CLIENT_ID"))
	values.Add("response_type", response_type)
	values.Add("scope", scope)
	values.Add("redirect_uri", redirect_uri)
	values.Add("state", state)
	values.Add("nonce", nonce)

	resp, err := http.Get(oauthURL + "?" + values.Encode())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}
