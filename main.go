package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
  "regexp"
)


func main() {

	r := mux.NewRouter()

	r.HandleFunc("/{domain}", ReadHandler).Methods("GET")
	r.HandleFunc("/{domain}/{addr}", WriteHandler).Methods("PUT")
	http.Handle("/", r)

	// Start the web server on port 8080
	port := 8080
	fmt.Printf("Server is listening on http://localhost:%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func ReadHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	domain := vars["domain"]

	if isValidDomain(domain) {
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, "Domain = %v", domain)
	} else {
		http.Error(resp, "Invalid domain format", http.StatusBadRequest)
	}
}

func WriteHandler(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	domain := vars["domain"]
	addr := vars["addr"]

	if isValidDomain(domain) && isValidAddr(addr) {
		resp.WriteHeader(http.StatusOK)
		fmt.Println("Domain = %v, Addr = %v", domain, addr)
		fmt.Fprintf(resp, "Domain = %v, Addr = %v", domain, addr)
	} else {
		http.Error(resp, "Invalid domain or address format", http.StatusBadRequest)
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
