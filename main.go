package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body
func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't, use
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the handler would keep
	// executing and also write the "Hello from SnippetBox" message
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the  value of the id parameter from the query string and try to convert it to an interger using the strconv.Atoi() function
	// If it can't be converted to an integer, or the value is less than 1, we return a 404 page not found response
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// NotFound replies to the request with an HTTP 404 not found error
		http.NotFound(w, r)
		return
	}

	// Use the fmt.Fprintf() function to interpolate the id value with our response and write it to the http.ResponseWriter
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not. Note that
	// http.MethodPost is a constant equal to the string "POST"
	if r.Method != http.MethodPost {
		// If it's not, use the w.WriteHeader() method to send a 405 status code
		// and the w.Write() method to write a "Method Not Allowed" response body. We then return from the function so that the subsequent code is not executed
		// Use the Header().Set() method to add an 'Allow: POST' header to the response header map. The first parameter is the header name, and the second parameter is the header value.
		// Use the http.Error() function to send a 405 status code and "Method Not Allowed" string as the response body
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		return
	}

	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the handler functions
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// User the http.ListenAndServe() function to start a web server. We pass in
	// two parameters: the TCP network address to listen on (in this case: ":4000")
	// and the servermux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
