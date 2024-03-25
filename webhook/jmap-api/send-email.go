package jmap_api

import (
	"fmt"
	"git.sr.ht/~rockorager/go-jmap"
	_ "git.sr.ht/~rockorager/go-jmap/core"
	"git.sr.ht/~rockorager/go-jmap/mail"
	"git.sr.ht/~rockorager/go-jmap/mail/email"
	"git.sr.ht/~rockorager/go-jmap/mail/emailsubmission"
)

const sessionEndpoint = "http://james.appscode.ninja:80/jmap/session"
const bearerToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0dXNlci5vcmdAbXlkb21haW4iLCJpYXQiOjE2MTYyMzkwMjIsImV4cCI6MTk0NjIzOTAyMn0.SO2fTuxB6qQNZqkhwWAndyqXDpyZMrHZFNdgG5ka7j8dW0G9JqU-rx04wqL-R2T51QCVYnlLBM62AVEltK-85ck10nFHEEFwQSD85v5PpI_qzeTz4NRFDxek-5N2gYdK0XQ6NzMvlSwxyUM2TkNZCupIto9H3MDeuDSs0p6Xmwk3iaGKJcIorIUbImf8xzwfB39ytpOw1j6ggAGjczZH8Ykz2PnHQQX2TU8R_dGbn6euYOqnTiZDWggfbxaS1joJ3PatN0Q-jxfuQGwTdZWeSN-Ocvr55MitQvMKJQaoRMVAOFCPXuhCkMq1szKRaBbhnL1nFjmHBpMJC5VpSeDgKOa-govmuTrnRMDV15n5KeeSWEeE2Km9ibOqzZktlR5EU-lU16h-a0u5ydPfax4HcUnNnVcKFqjSGyMFyBQOOxQWKrSJdHNaA7ZP07QQjwdSRQFbJHsUBdH22oBfTnT891yMmWW1iFLHZuV2sMivUGnH-EBO29HHVzAXUlOzErqDaVnovMvdPN6_Vi80LpdSWikYbe2zXUV4csBAc6bC1NweP9omNidBk9Vgo-3q2mHw-H6mfFnTONUipbtbV8ifI2MKQ-ZS0ciNomOcNgyYu5mc-ebDFpZU9xx-_nakCOscUm0tHN224MxKNIx-9E7XrS0PlSOauNYqi4AVKP9svW0"
const draftMailboxID = "288b8a70-e82d-11ee-8288-23209b563714"
const sentMailboxID = "28808df0-e82d-11ee-8288-23209b563714"
const userEmail = "testuser.org@mydomain"

var Client *jmap.Client
var UserID jmap.ID

func init() {
	Client = &jmap.Client{
		SessionEndpoint: sessionEndpoint,
	}

	Client.WithAccessToken(bearerToken)

	if err := Client.Authenticate(); err != nil {
		panic(err)
	}

	UserID = Client.Session.PrimaryAccounts[mail.URI]
}

func SendEmail(myMail *email.Email) {
	req := &jmap.Request{
		Using: []jmap.URI{"urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"},
	}

	invokeSetDraftEMail(req, UserID, myMail)
	invokeSendEmail(req, UserID)

	resp, err := Client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.CreatedIDs)
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
				"mailboxIds/" + draftMailboxID: nil,
				"mailboxIds/" + sentMailboxID:  true,
				"keywords/$seen":               true,
				"keywords/$draft":              nil,
			},
		},
	})
}
