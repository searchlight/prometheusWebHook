package jmap_api

import (
	"git.sr.ht/~rockorager/go-jmap"
	_ "git.sr.ht/~rockorager/go-jmap/core"
	"git.sr.ht/~rockorager/go-jmap/mail"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
	"git.sr.ht/~rockorager/go-jmap/mail/emailsubmission"
	"log"
)

const sessionEndpoint = "http://james.appscode.ninja:80/jmap/session"
const bearerToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0dXNlci5vcmdAbXlkb21haW4iLCJpYXQiOjE2MTYyMzkwMjIsImV4cCI6MTk0NjIzOTAyMn0.SO2fTuxB6qQNZqkhwWAndyqXDpyZMrHZFNdgG5ka7j8dW0G9JqU-rx04wqL-R2T51QCVYnlLBM62AVEltK-85ck10nFHEEFwQSD85v5PpI_qzeTz4NRFDxek-5N2gYdK0XQ6NzMvlSwxyUM2TkNZCupIto9H3MDeuDSs0p6Xmwk3iaGKJcIorIUbImf8xzwfB39ytpOw1j6ggAGjczZH8Ykz2PnHQQX2TU8R_dGbn6euYOqnTiZDWggfbxaS1joJ3PatN0Q-jxfuQGwTdZWeSN-Ocvr55MitQvMKJQaoRMVAOFCPXuhCkMq1szKRaBbhnL1nFjmHBpMJC5VpSeDgKOa-govmuTrnRMDV15n5KeeSWEeE2Km9ibOqzZktlR5EU-lU16h-a0u5ydPfax4HcUnNnVcKFqjSGyMFyBQOOxQWKrSJdHNaA7ZP07QQjwdSRQFbJHsUBdH22oBfTnT891yMmWW1iFLHZuV2sMivUGnH-EBO29HHVzAXUlOzErqDaVnovMvdPN6_Vi80LpdSWikYbe2zXUV4csBAc6bC1NweP9omNidBk9Vgo-3q2mHw-H6mfFnTONUipbtbV8ifI2MKQ-ZS0ciNomOcNgyYu5mc-ebDFpZU9xx-_nakCOscUm0tHN224MxKNIx-9E7XrS0PlSOauNYqi4AVKP9svW0"

var draftMailboxID jmap.ID
var sentMailboxID jmap.ID
var userEmail string

var myClient *jmap.Client
var userID jmap.ID

func init() {
	myClient = &jmap.Client{
		SessionEndpoint: sessionEndpoint,
	}

	myClient.WithAccessToken(bearerToken)

	if err := myClient.Authenticate(); err != nil {
		log.Fatal("unable to authenticate user with the given credentials", err)
	}

	userID = myClient.Session.PrimaryAccounts[mail.URI]
	userEmail = myClient.Session.Accounts[userID].Name

	///Not using RFC8621's "Mailbox/query filter" because it's not properly implemented in the James distributed server 3.8.0
	var err error
	if draftMailboxID, err = getMailboxIdByTag("Drafts"); err != nil {
		log.Fatal(err)
	}

	if sentMailboxID, err = getMailboxIdByTag("Sent"); err != nil {
		log.Fatal(err)
	}
}

func SendEmail(myMail *email.Email) error {
	req := &jmap.Request{
		Using: []jmap.URI{"urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"},
	}

	invokeSetDraftEMail(req, userID, myMail)
	invokeSendEmail(req, userID)

	if _, err := myClient.Do(req); err != nil {
		return err
	}

	return nil
}

func invokeSetDraftEMail(req *jmap.Request, id jmap.ID, myMail *email.Email) {
	myMap := map[jmap.ID]*email.Email{
		"draft": &(*myMail),
	}

	req.Invoke(&email.Set{
		Account: id,
		Create:  myMap,
	})
}

func invokeSendEmail(req *jmap.Request, id jmap.ID) {
	myEmailSubmission := emailsubmission.EmailSubmission{
		EmailID: "#draft",
	}

	req.Invoke(&emailsubmission.Set{
		Account: id,

		Create: map[jmap.ID]*emailsubmission.EmailSubmission{
			"sendIt": &myEmailSubmission,
		},

		OnSuccessUpdateEmail: map[jmap.ID]jmap.Patch{
			"#sendIt": {
				"mailboxIds/" + string(draftMailboxID): nil,
				"mailboxIds/" + string(sentMailboxID):  true,
				"keywords/$seen":                       nil,
				"keywords/$draft":                      nil,
			},
		},
	})
}
