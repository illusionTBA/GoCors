package main

import (
	"io"
	"log"
	"net/http"
	"strings"
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

		// Get the destination URL from the request path
		destinationURL := r.URL.Path[1:] // Remove the leading forward slash
		log.Printf(destinationURL)
		// Get the headers from the query parameters
		headers := r.URL.Query().Get("headers")

		// Forward the request to the desired destination URL
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

		// Add user-supplied headers to the destination request
		if headers != "" {
			headerList := strings.Split(headers, ",")
			for _, header := range headerList {
				headerParts := strings.SplitN(header, ":", 2)
				if len(headerParts) == 2 {
					key := strings.TrimSpace(headerParts[0])
					value := strings.TrimSpace(headerParts[1])
					req.Header.Set(key, value)
				}
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
			// Exclude duplicate CORS headers
			if !strings.HasPrefix(strings.ToLower(key), "access-control-") {
				for _, value := range values {
					w.Header().Add(key, value)
				}
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

	// Start the server on port 3000
	log.Printf("Server up @ localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}
