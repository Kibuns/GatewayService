package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the query parameter from the URL
	vars := mux.Vars(r)
	query := vars["query"]

	// Make a GET request to Service2 with the query parameter
	resp, err := http.Get("http://localhost:9999/search/" + query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers from Service2 to the gateway response
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}

	// Copy the status code from Service2 to the gateway response
	w.WriteHeader(resp.StatusCode)

	// Copy the response body from Service2 to the gateway response
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}

// Define a handler function for the /service1 endpoint
func searchHomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("search home")
	// Make a request to Service 1
	resp, err := http.Get("http://localhost:9999")
	if err != nil {
		http.Error(w, "Error calling Service 1", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers from Service 1 to the gateway response
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}

	// Copy the status code from Service 1 to the gateway response
	w.WriteHeader(resp.StatusCode)

	// Copy the response body from Service 1 to the gateway response
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}

func twootHomeHandler(w http.ResponseWriter, r *http.Request) {
	// Make a request to Service 1
	resp, err := http.Get("http://localhost:10000")
	if err != nil {
		http.Error(w, "Error calling Service 1", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers from Service 1 to the gateway response
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}

	// Copy the status code from Service 1 to the gateway response
	w.WriteHeader(resp.StatusCode)

	// Copy the response body from Service 1 to the gateway response
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}

func storeTwootHandler(w http.ResponseWriter, r *http.Request) {
	// Make a POST request to Service1 with the request body
	resp, err := http.Post("http://localhost:10000/create", "application/json", r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers from Service1 to the gateway response
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}

	// Copy the status code from Service1 to the gateway response
	w.WriteHeader(resp.StatusCode)

	// Copy the response body from Service1 to the gateway response
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}

func main() {
    // Create a new HTTP server
    server := mux.NewRouter().StrictSlash(true)

	server.Use(CORS)

    // Route requests to the appropriate handler function based on the URL path
    server.HandleFunc("/search", searchHomeHandler)
	server.HandleFunc("/search/{query}", searchHandler)
    server.HandleFunc("/twoot", twootHomeHandler)
	server.HandleFunc("/twoot/store", storeTwootHandler)


    // Start the HTTP server on port 8080
    log.Println("API Gateway listening on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", server))

	
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		//return
	})

}
