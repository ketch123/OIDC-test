package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const (
	rs2Letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	oauthURL      = "https://accounts.google.com/o/oauth2/v2/auth"
	response_type = "code"
	scope         = "openid email profile"
	redirect_uri  = "http://localhost:8080/tokenReq"
)

var state string = RandString(15)
var nonce string = RandString(20)
var client_id = os.Getenv("CLIENT_ID")
var client_secret = os.Getenv("CLIENT_SECRET")

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", authReq)
	router.Handle("/", router)
	router.HandleFunc("/tokenReq", tokenReq)
	router.Handle("/tokenReq", router)
	router.HandleFunc("/accessTokenReq", accessTokenReq)
	router.Handle("/accessTokenReq", router)
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

func authReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("authReq started...")

	values := url.Values{}
	values.Add("client_id", client_id)
	values.Add("response_type", response_type)
	values.Add("scope", scope)
	values.Add("redirect_uri", redirect_uri)
	values.Add("state", state)
	values.Add("nonce", nonce)

	http.Redirect(w, r, oauthURL+"?"+values.Encode(), http.StatusFound)
}

func tokenReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tokenReq started ...")

	if r.URL.Query().Get("state") != state {
		fmt.Println("state didn't match")
	}

	client := &http.Client{Timeout: time.Duration(10) * time.Second}
	host_uri := "https://www.googleapis.com/oauth2/v4/token"
	code := r.URL.Query().Get("code")
	grant_type := "authorization_code"

	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", client_id)
	values.Add("client_secret", client_secret)
	values.Add("redirect_uri", redirect_uri)
	fmt.Printf("redirect_uri = %s\n", redirect_uri)
	values.Add("grant_type", grant_type)

	req, err := http.NewRequest("POST", host_uri, strings.NewReader(values.Encode()))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func accessTokenReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("access starde...")
}
