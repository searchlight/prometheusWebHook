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

	fmt.Println("Alert received")

	//myPodLogs := string(podLogs.GetPodLogs("default")[:])
	myPodLogs, err := podLogs.GetPodLogs("default")
	if err != nil {
		http.Error(w, "Error fetching pod logs", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Read the request body
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r.Body); err != nil {
		http.Error(w, "Alert body missing/invalid", http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}
	body := buf.String()

	myMail, err := jmap_api.NewEmailBuilder().
		SetSubject("Prometheus Alertmanager alert received").
		SetBodyValue(body).
		SetAttachment(bytes.NewReader(myPodLogs)).
		SetRecipient("testuser.org@mydomain").
		Build()

	if err != nil {
		http.Error(w, "Unable to create email", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(myMail)
	if err := jmap_api.SendEmail(&myMail); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		fmt.Println("Failed to send email", err)
		return
	}
	// Respond with HTTP status 200 OK
	w.WriteHeader(http.StatusOK)
	fmt.Println("Alert email sent successfully")
}
