package jmap_api

import (
	"errors"
	"fmt"
	"git.sr.ht/~rockorager/go-jmap"
	_ "git.sr.ht/~rockorager/go-jmap/core"
	"git.sr.ht/~rockorager/go-jmap/mail"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"strings"
)

type emailBuilder struct {
	recipient, subject, bodyValue string
	uploadResponse                *jmap.UploadResponse
}

func NewEmailBuilder() *emailBuilder {
	return &emailBuilder{}
}

func (b *emailBuilder) SetSubject(subject string) *emailBuilder {
	b.subject = subject
	return b
}

func (b *emailBuilder) SetBodyValue(body string) *emailBuilder {
	b.bodyValue = body
	return b
}

func (b *emailBuilder) SetRecipient(recipient string) *emailBuilder {
	b.recipient = recipient
	return b
}

func (b *emailBuilder) SetAttachment(blob io.Reader) *emailBuilder {
	resp, err := Upload(myClient, userID, blob)
	b.uploadResponse = nil

	if err != nil {
		fmt.Println("Error setting attachment ", err.Error())
		return b
	}

	if resp == nil {
		fmt.Println("response is nil")
		return b
	}

	b.uploadResponse = resp
	return b
}

// /Slightly modified version of https://github.com/rockorager/go-jmap/blob/main/client.go#L162
func Upload(c *jmap.Client, accountID jmap.ID, blob io.Reader) (*jmap.UploadResponse, error) {
	c.Lock()
	if c.SessionEndpoint == "" {
		c.Unlock()
		return nil, fmt.Errorf("jmap/client: SessionEndpoint is empty")
	}
	if c.Session == nil {
		c.Unlock()
		err := c.Authenticate()
		if err != nil {
			return nil, err
		}
	}

	url := strings.ReplaceAll(c.Session.UploadURL, "{accountId}", string(accountID))
	c.Unlock()
	req, err := http.NewRequest("POST", url, blob)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode == 200 || resp.StatusCode == 201) {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	info := &jmap.UploadResponse{}
	err = json.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (b *emailBuilder) Build() (email.Email, error) {
	if b.recipient == "" {
		return email.Email{}, errors.New("No recipient defined")
	}

	if b.uploadResponse == nil {
		return email.Email{}, errors.New("Error setting attachment")
	}

	from := mail.Address{
		Name:  userEmail,
		Email: userEmail,
	}

	to := mail.Address{
		Name:  b.recipient,
		Email: b.recipient,
	}

	myBodyValue := email.BodyValue{
		Value: b.bodyValue,
	}

	myBodyPart := email.BodyPart{
		PartID: "body",
		Type:   "text/plain",
	}

	myAttachment := email.BodyPart{
		BlobID: b.uploadResponse.ID,

		Size: b.uploadResponse.Size,

		Type: b.uploadResponse.Type,

		Name: "pod_Logs.txt",

		Disposition: "attachment",
	}

	//_ = myAttachment

	myMail := email.Email{
		From: []*mail.Address{
			&from,
		},

		To: []*mail.Address{
			&to,
		},

		Subject: b.subject,

		Keywords: map[string]bool{"$draft": true},

		MailboxIDs: map[jmap.ID]bool{jmap.ID(draftMailboxID): true},

		BodyValues: map[string]*email.BodyValue{"body": &myBodyValue},

		TextBody: []*email.BodyPart{&myBodyPart},

		HasAttachment: true,

		Attachments: []*email.BodyPart{&myAttachment},
	}

	return myMail, nil
}
