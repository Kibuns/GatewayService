package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Kibuns/GatewayService/Models"
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
	body := r.Body
	fmt.Println("Storing Twoot")

	// Parse the request body into a Twoot struct
	var twoot Models.Twoot
	err := json.NewDecoder(body).Decode(&twoot)
	fmt.Println(twoot.Content)
	if err != nil {
		http.Error(w, "Could not decode body into twoot", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	// Clone the incoming request
	newReq, err := http.NewRequest("GET", "http://localhost:3500/getusername", nil)
	if err != nil {
		http.Error(w, "Failed to create new request", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Set the headers of the new request to match the incoming request
	newReq.Header = r.Header.Clone()

	// Make the request to retrieve the username
	client := http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//if status code is not okay, means error message is in the body
		valueBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, "Failed to read username response", http.StatusInternalServerError)
				fmt.Println(err)
				return
			}

		value := string(valueBytes)
		http.Error(w, value, http.StatusInternalServerError)
		fmt.Println("Failed to fetch username")
		return
	}

	// Read the response body
	usernameBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read username response", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	username := string(usernameBytes)
	twoot.UserName = username

	// Convert the twoot object to JSON
	payload, err := json.Marshal(twoot)
	if err != nil {
		http.Error(w, "Failed to encode twoot payload", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Make a POST request to Service1 with the twoot payload
	resp, err = http.Post("http://localhost:10000/create", "application/json", bytes.NewBuffer(payload))
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

func createUserHandler(w http.ResponseWriter, r *http.Request){
		// Make a POST request to Service1 with the request body
		resp, err := http.Post("http://localhost:9998/create", "application/json", r.Body)
		if err != nil {
			http.Error(w, "Could not create user", http.StatusInternalServerError)
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

func getJWTHandler(w http.ResponseWriter, r *http.Request) {
	// Make a POST request to Service1 with the request body
	resp, err := http.Post("http://localhost:3500/jwt", "application/json", r.Body)
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

	// server.Use(CORS)

    // Route requests to the appropriate handler function based on the URL path
    server.HandleFunc("/search", searchHomeHandler)
	server.HandleFunc("/search/{query}", searchHandler)
    server.HandleFunc("/twoot", twootHomeHandler)
	server.HandleFunc("/twoot/post", storeTwootHandler)
	server.HandleFunc("/user/create", createUserHandler)
	server.HandleFunc("/jwt", getJWTHandler)


    // Start the HTTP server on port 8080
    log.Println("API Gateway listening on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", server))

	
}

// func CORS(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		// Set headers
// 		w.Header().Set("Access-Control-Allow-Headers:", "*")
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "*")

// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		fmt.Println("ok")

// 		// Next
// 		next.ServeHTTP(w, r)
// 		//return
// 	})

// }
