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

const (
	twootServiceURL = "http://twootservice:10000";
	searchServiceURL = "http://searchservice:8081";
	userServiceURL = "http://userservice:9998";
	authServiceURL = "http://authservice:8083";
)


func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage of the Gateway Service!")
	fmt.Println("Endpoint Hit: gateway")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the query parameter from the URL
	vars := mux.Vars(r)
	query := vars["query"]

	// Make a GET request to Service2 with the query parameter
	resp, err := http.Get(searchServiceURL + "/search/" + query)
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
	resp, err := http.Get(searchServiceURL)
	if err != nil {
		http.Error(w, "Error calling search service", http.StatusInternalServerError)
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
	// resp, err := http.Get("http://localhost:10000")
	resp, err := http.Get(twootServiceURL)
	if err != nil {
		http.Error(w, "Error calling twoot service", http.StatusInternalServerError)
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
	newReq, err := http.NewRequest("GET", authServiceURL + "/getusername", nil)
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
	resp, err = http.Post(twootServiceURL + "/create", "application/json", bytes.NewBuffer(payload))
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

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Define a struct to unmarshal the request body
	type UserRequest struct {
		PermissionToSave bool `json:"permissionToSave"`
		Username string `json:"username"`
		Password string `json:"password"`
	}


	var userReq UserRequest
	if err := json.Unmarshal(body, &userReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check the value of "permissionToSave" field
	if userReq.PermissionToSave {
		// The value of "permissionToSave" is true
		// Make the POST request to Service1
		resp, err := http.Post(userServiceURL+"/create", "application/json", bytes.NewReader(body))
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
	} else {
		// The value of "permissionToSave" is not true
		// Handle the case when permission is not granted
		http.Error(w, "Permission not granted to save", http.StatusForbidden)
	}
}

func getJWTHandler(w http.ResponseWriter, r *http.Request) {
	// Make a POST request to Service1 with the request body
	resp, err := http.Post(authServiceURL + "/jwt", "application/json", r.Body)
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

func deleteAllOfUser(w http.ResponseWriter, r *http.Request) {
	// Extract the query parameter from the URL
	vars := mux.Vars(r)
	username := vars["username"]


    // Publish the event to RabbitMQ
    send(username)

	fmt.Fprintf(w, "Sent message to delete user: " + username)
}

func getAllOfUser(w http.ResponseWriter, r *http.Request) {
	// Extract the query parameter from the URL
	vars := mux.Vars(r)
	username := vars["username"]

	twootResp, err := http.Get(twootServiceURL + "/getall/" + username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer twootResp.Body.Close()

	userResp, err := http.Get(userServiceURL + "/get/" + username)
	if err != nil {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	var twoots []Models.Twoot
	var user Models.User

	err = json.NewDecoder(twootResp.Body).Decode(&twoots)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(userResp.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type CombinedResponse struct {
		Twoots []Models.Twoot `json:"twoots"`
		User   Models.User    `json:"user"`
	}

	combinedResponse := CombinedResponse{
		Twoots: twoots,
		User:   user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(combinedResponse)
}

func main() {
    // Create a new HTTP server
    server := mux.NewRouter().StrictSlash(true)

	// server.Use(CORS)

    // Route requests to the appropriate handler function based on the URL path
	server.HandleFunc("/", homePage)
    server.HandleFunc("/search", searchHomeHandler)
	server.HandleFunc("/search/{query}", searchHandler)
	server.HandleFunc("/delete/{username}", deleteAllOfUser)
	server.HandleFunc("/getall/{username}", getAllOfUser)
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
