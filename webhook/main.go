package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	podLogs "webhook/client-go"
	jmap_api "webhook/jmap-api"
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
	buf := new(strings.Builder)
	io.Copy(buf, r.Body)
	body := buf.String()

	// Print the request body
	fmt.Println("Received webhook payload:")
	fmt.Println(body)

	fmt.Println("Here is the logs\n\n\n")

	//myPodLogs := string(podLogs.GetPodLogs("default")[:])
	myPodLogs := podLogs.GetPodLogs("default")

	myMail := jmap_api.NewEmailBuilder().
		SetSubject("Prometheus Alertmanager alert received").
		SetBodyValue(body).
		SetAttachment(bytes.NewReader(myPodLogs)).
		SetRecipient("testuser.org@mydomain").
		Build()

	//fmt.Println(myMail)
	jmap_api.SendEmail(&myMail)
	// Respond with HTTP status 200 OK
	w.WriteHeader(http.StatusOK)
}
