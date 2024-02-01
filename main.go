package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	Port    int
	BaseUrl string
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

func init() {
	BaseUrl = getEnv("BASE_URL", "http://127.19.73.21:17468")
	Port = getEnvAsInt("PORT", 8080)
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
	// validate
	vars := mux.Vars(r)
	domain := vars["domain"]
	if !isValidDomain(domain) {
		http.Error(w, "invalid format", http.StatusBadRequest)
		return
	}
	// send
	url := fmt.Sprintf("%s/state/%s", BaseUrl, domain)
	fmt.Printf("Sending request: GET %v...", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	handle(resp, w)
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	// validate inputs
	vars := mux.Vars(r)
	domain := vars["domain"]
	addr := vars["addr"]
	if !isValidDomain(domain) || !isValidAddr(addr) {
		http.Error(w, "Invalid format.", http.StatusBadRequest)
		return
	}

	// prepare request
	jsonData, err := getPayload(domain, addr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tokenData, err := getToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("%s/tx", BaseUrl)
	fmt.Printf("Sending request: PUT %v...", url)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Header.Set(tokenData.Header, tokenData.Token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	// send request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Sending request failed", http.StatusBadRequest)
	}
	defer resp.Body.Close()
	handle(resp, w)
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
		http.Error(w, fmt.Sprintf("Server response status: %v.", response.StatusCode), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if len(body) > 0 {
		fmt.Fprint(w, string(body))
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

func getToken() (TokenData, error) {
	var tokenData TokenData
	url := fmt.Sprintf("%s/testnet/token", BaseUrl)
	fmt.Printf("Sending request: GET %v...", url)
	tokenResp, err := http.Get(url)
	if err != nil {
		return tokenData, err
	}
	body, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		return tokenData, err
	}
	err = json.Unmarshal([]byte(body), &tokenData)
	if err != nil {
		return tokenData, err
	}
	return tokenData, nil
}

func getPayload(domain string, addr string) ([]byte, error) {
	payload := Payload{
		Command: "data",
		NS:      domain,
		D:       addr,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}
