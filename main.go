package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// Define the handler function that will handle the CORS requests
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		// If it's a preflight request, respond with 200 OK
		if r.Method == "OPTIONS" {
			return
		}

		// Forward the request to the desired destination URL
		destinationURL := r.URL.Path[1:] // Remove the leading forward slash
		req, err := http.NewRequest(r.Method, destinationURL, r.Body)
		if err != nil {
			log.Printf("Error creating request: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Copy headers from the original request to the destination request
		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Send the request to the destination server
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error forwarding request: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy the response headers from the destination server to the response writer
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Set the status code of the response
		w.WriteHeader(resp.StatusCode)

		// Copy the response body from the destination server to the response writer
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("Error writing response: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":3000", handler))
}
