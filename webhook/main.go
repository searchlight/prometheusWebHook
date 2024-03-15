package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	podLogs "webhook/client-go"
)

func main() {
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Print the request body
	fmt.Println("Received webhook payload:")
	fmt.Println(string(body))

	fmt.Println("Here is the logs\n\n\n")

	myPodLogs := podLogs.GetPodLogs("default")
	for i, _ := range myPodLogs {
		fmt.Println(myPodLogs[i])
	}

	// Respond with HTTP status 200 OK
	w.WriteHeader(http.StatusOK)
}
