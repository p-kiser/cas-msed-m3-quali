package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

const (
	Port    = 8080
	BaseUrl = "http://127.19.73.21:17468" // TODO: check this out: http://127.19.73.21:17468/about
)

type TokenData struct {
	Header string `json:"header"`
	Token  string `json:"token"`
}

type Payload struct {
	Command string `json:"command"`
	NS      string `json:"ns"`
	D       string `json:"d"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/{domain}", GetHandler).Methods("GET")
	router.HandleFunc("/{domain}/{addr}", PutHandler).Methods("PUT")
	http.Handle("/", router)

	fmt.Printf("Server is listening on http://localhost:%d...\n", Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", Port), nil)

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := vars["domain"]
	if !isValidDomain(domain) {
		http.Error(w, "invalid domain", http.StatusBadRequest)
		return
	}
	url := fmt.Sprintf("%s/state/%s", BaseUrl, domain)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("GET request failed: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	handle(resp, w)
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := vars["domain"]
	addr := vars["addr"]

	if !isValidDomain(domain) || !isValidAddr(addr) {
		http.Error(w, "Invalid domain or address format", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/tx", BaseUrl)

	payload := Payload{
		Command: "data",
		NS:      domain,
		D:       addr,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	dataBytes := []byte(string(jsonData))

	fmt.Println(string(jsonData))

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(dataBytes))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	// TODO: get token /testnet/token
	tokenUrl := fmt.Sprintf("%s/testnet/token", BaseUrl)
	tokenResp, err := http.Get(tokenUrl)
	if err != nil {
		http.Error(w, fmt.Sprintf("token request failed: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	body, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		fmt.Println("Error reading token response body:", err)
		return
	}

	var tokenData TokenData

	// Parse the JSON into the struct
	err = json.Unmarshal([]byte(body), &tokenData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	req.Header.Set(tokenData.Header, tokenData.Token)
	req.Header.Set("Content-Type", "application/json")


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Sending request failed", http.StatusBadRequest)
	}
	defer resp.Body.Close()
	handle(resp, w)
	return
}

func handle(response *http.Response, w http.ResponseWriter) {

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Printf("API Response (%d):\n", response.StatusCode)
	if len(body) > 0 {
		fmt.Println("---")
		fmt.Println(string(body))
		fmt.Println("---")
	}

	if response.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Server error: %v.", response.StatusCode), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(body))
	}
}

func isValidDomain(domain string) bool {
	pattern := "^[a-z0-9-_]{3,64}\\.i2p$"
	return isValid(pattern, domain)
}

func isValidAddr(addr string) bool {
	pattern := "^[a-z0-9]{52}$"
	return isValid(pattern, addr)
}

func isValid(pattern string, input string) bool {
	regexpPattern := regexp.MustCompile(pattern)
	return regexpPattern.MatchString(input)
}
