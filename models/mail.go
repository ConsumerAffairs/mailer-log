package models

import (
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"time"
)

type MailAttachment struct {
	Filetype string `json:"filetype"`
	Mimetype string `json:"mimetype"`
}

type Mail struct {
	Id             bson.ObjectId    `json:"id"              bson:"_id,omitempty"`
	Sent_at        time.Time        `json:"sentAt"`
	From_email     string           `json:"fromEmail"`
	To_emails      []string         `json:"toEmails"`
	Cc_emails      []string         `json:"ccEmails"`
	Bcc_emails     []string         `json:"bccEmails"`
	All_recipients []string         `json:"allRecipients"`
	Headers        string           `json:"headers"`
	Subject        string           `json:"subject"`
	Body_html      template.HTML    `json:"bodyHtml"`
	Body_text      string           `json:"bodyText"`
	Raw            string           `json:"raw"`
	Sent           bool             `json:"sent"`
	Attachaments   []MailAttachment `json:"attachaments"`
}
