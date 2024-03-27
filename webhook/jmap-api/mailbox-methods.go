package jmap_api

import (
	"errors"
	"fmt"
	"git.sr.ht/~rockorager/go-jmap"
	"git.sr.ht/~rockorager/go-jmap/mail/mailbox"
)

// ##### UNTESTED #####
func getAllMailboxes() ([]*mailbox.Mailbox, error) {
	req := &jmap.Request{
		Using: []jmap.URI{"urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"},
	}

	req.Invoke(&mailbox.Get{
		Account: userID,
	})

	resp, err := myClient.Do(req)
	if err != nil || resp == nil {
		fmt.Println("Error fetching mailboxes")
		return []*mailbox.Mailbox{}, err
	}

	if len(resp.Responses) > 1 {
		return []*mailbox.Mailbox{}, errors.New("Multiple responses received on function getAllMailboxes")
	}

	///len(resp.Responses) == 1
	var mailboxList []*mailbox.Mailbox

	for _, inv := range resp.Responses {
		switch r := inv.Args.(type) {
		case *mailbox.GetResponse:
			mailboxList = r.List
			break
		}
	}

	return mailboxList, nil
}

func getMailboxIdByTag(tag string) (jmap.ID, error) {
	req := &jmap.Request{
		Using: []jmap.URI{"urn:ietf:params:jmap:core", "urn:ietf:params:jmap:mail"},
	}

	myFilter := mailbox.FilterCondition{
		Role: mailbox.Role(tag),
	}

	req.Invoke(&mailbox.Query{
		Account: userID,
		Filter:  &myFilter,
	})

	resp, err := myClient.Do(req)
	if err != nil || resp == nil {
		fmt.Println("error fetching mailbox with tag: " + tag)
		return "", err
	}

	if len(resp.Responses) > 1 {
		return "", errors.New("Multiple responses received on function getMailboxIdByTag with tag: " + tag)
	}

	///len(resp.Responses) == 1
	var requiredID jmap.ID

	for _, inv := range resp.Responses {
		switch r := inv.Args.(type) {
		case *mailbox.QueryResponse:
			if len(r.IDs) > 1 {
				return "", errors.New("Multiple mailboxes exists with the same role: " + tag)
			}
			requiredID = r.IDs[0]
			break
		}
	}

	return requiredID, nil
}
