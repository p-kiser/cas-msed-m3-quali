package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/{domain}", ReadHandler).Methods("GET")
	// r.HandleFunc("/{domain:[a-z0-9-_]{3,64}\\.i2p}", ReadHandler)

	// r.HandleFunc("/{domain:[a-z0-9-_{3-64}]\.i2p}/{addr:[a-z0-9]{52}}", WriteHandler)
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
	domain _= vars["domain"]

	if isValidDomain() {
		fmt.Fprintf(resp, "Domain = %v", domain)
		resp.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid domain format", http.StatusBadRequest)
	}
}

func isValidDomain(domain string) bool {
	pattern := "^[a-z0-9-_]{3,64}\\.i2p$"
	regexpPattern := regexp.MustCompile(pattern)
	return regexpPattern.MatchString(domain)
}
