package main

import (
	"encoding/base64"
	"encoding/json"
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

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
}

type Id_token struct {
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash""`
	Nonce         string `json:"nonce"`
	Exp           int    `json:"exp"`
	Iss           string `json:"iss"`
	Iat           int    `json:"iat"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
}

const (
	rs2Letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	authURI      = "https://accounts.google.com/o/oauth2/v2/auth"
	responseType = "code"
	scope        = "openid email profile"
	redirectURI  = "http://localhost:8080/tokenReq"
)

var state string = RandString(15)
var nonce string = RandString(20)
var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")

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
	values.Add("client_id", clientID)
	values.Add("response_type", responseType)
	values.Add("scope", scope)
	values.Add("redirect_uri", redirectURI)
	values.Add("state", state)
	values.Add("nonce", nonce)

	http.Redirect(w, r, authURI+"?"+values.Encode(), http.StatusFound)
}

func tokenReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tokenReq started ...")

	if r.URL.Query().Get("state") != state {
		log.Fatal("state did not match")
	}

	client := &http.Client{Timeout: time.Duration(10) * time.Second}
	tokenURI := "https://www.googleapis.com/oauth2/v4/token"
	code := r.URL.Query().Get("code")
	grantType := "authorization_code"

	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", clientID)
	values.Add("client_secret", clientSecret)
	values.Add("redirect_uri", redirectURI)
	values.Add("grant_type", grantType)

	req, err := http.NewRequest("POST", tokenURI, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	getResource(body)
}

func getResource(body []uint8) {
	fmt.Println("getResource started...")
	//構造体の初期化を行い、jsonを構造体に埋め込む
	tokens := Token{}
	if err := json.Unmarshal(body, &tokens); err != nil {
		log.Fatal(err)
	}

	parts := strings.Split(tokens.IdToken, ".")
	data, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Fatal(err)
	}

	idTokens := Id_token{}
	if err := json.Unmarshal(data, &idTokens); err != nil {
		log.Fatal(err)
	}
	if nonce != idTokens.Nonce {
		log.Fatal("nonce did not match")
	}

	fmt.Printf("Hello, %s!!\n", idTokens.Name)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", authReq)
	router.Handle("/", router)
	router.HandleFunc("/tokenReq", tokenReq)
	router.Handle("/tokenReq", router)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", "8080"), router); err != nil {
		log.Fatal(err)
	}
}
