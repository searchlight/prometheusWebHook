package jmap_api

import (
	"fmt"
	"git.sr.ht/~rockorager/go-jmap"
	_ "git.sr.ht/~rockorager/go-jmap/core"
	"git.sr.ht/~rockorager/go-jmap/mail"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
	_ "git.sr.ht/~rockorager/go-jmap/mail/emailsubmission"
	"io"
)

type emailBuilder struct {
	recipient, subject, bodyValue string
	uploadResponse                jmap.UploadResponse
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
	/*{
		buf := new(strings.Builder)
		io.Copy(buf, blob)
		fmt.Println(buf.String())
	}*/
	resp, err := Client.Upload(UserID, blob)
	//fmt.Println(resp)
	fmt.Println(err)
	if err != nil {
		fmt.Println("Error setting attachment ")
		//log.Fatal(err.Error())
		//fmt.Println("ERROR END ----------------------------------------")
	}

	if resp == nil {
		fmt.Println("response is nil")
	}
	_ = resp
	//b.uploadResponse = *resp
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

	myAttachment := email.BodyPart{
		BlobID: b.uploadResponse.ID,

		Size: b.uploadResponse.Size,

		Type: b.uploadResponse.Type,

		Name: "pod_Logs.txt",

		Disposition: "attachment",
	}

	_ = myAttachment

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

		//HasAttachment: true,

		//Attachments: []*email.BodyPart{&myAttachment},
	}

	return myMail
}
