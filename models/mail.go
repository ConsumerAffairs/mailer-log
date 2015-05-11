package models

import (
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"time"
)

type Mail struct {
	Id             bson.ObjectId `json:"id"              bson:"_id,omitempty"`
	Sent_at        time.Time     `json:"sent_at"`
	From_email     string        `json:"from_email"`
	To_emails      string        `json:"to_emails"`
	Cc_emails      string        `json:"cc_emails"`
	Bcc_emails     string        `json:"bcc_emails"`
	All_recipients string        `json:"all_recipients"`
	Headers        string        `json:"headers"`
	Subject        string        `json:"subject"`
	Body           template.HTML `json:"body"`
	Raw            string        `json:"raw"`
	Sent           string        `json:"sent"`
}

type MailResource struct {
	Data Mail `json:"data"`
}

type MailsCollection struct {
	Data []Mail `json:"data"`
}
