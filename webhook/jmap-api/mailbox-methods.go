package jmap_api

import (
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail/mailbox"
)

func getAllMailboxes() []*mailbox.Mailbox {
	req := &jmap.Request{
		Using: []jmap.URI{"urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"},
	}

	req.Invoke(&mailbox.Get{
		Account: userID,
	})

	//TODO: return a real list
	return []*mailbox.Mailbox{}
}
