package main

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Post("/webhook", webhookHandler)
	http.HandleFunc("/webhook", webhookHandler)
	fmt.Println("Webhook server started, listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	Body := buf.String()
	// Print the request body
	fmt.Println("Received webhook payload:")
	fmt.Println(string(Body))

	fmt.Println("-------------------------\nPRINTING POD LOGS: ")
	tmp := GetPodLogs("default")
	fmt.Println(len(tmp))
	for i, _ := range tmp {
		fmt.Println(tmp[i])
	}
	fmt.Println("-------------------------\n")
	// Respond with HTTP status 200 OK
	w.WriteHeader(http.StatusOK)
}
