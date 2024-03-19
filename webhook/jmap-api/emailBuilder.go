package jmap_api

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
)

type emailBuilder struct {
	recipient, subject, bodyValue string
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

func (b *emailBuilder) Build() email.Email {
	if b.recipient == "" {
		panic("no recipient defined")
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
	}

	return myMail
}
