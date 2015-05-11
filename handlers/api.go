package handlers

import (
	"encoding/json"
	"github.com/ConsumerAffairs/mailer_microservice/models"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type (
	MailController struct {
		session *mgo.Session
	}
)

// Return controllers with a mongodb session
func NewMailController(s *mgo.Session) *MailController {
	return &MailController{s}
}

//
// Mail handlers
//

func (uc MailController) ListMail(w http.ResponseWriter, r *http.Request) {
	mailsCollection := models.MailsCollection{[]models.Mail{}}

	err := uc.session.DB("mailer").C("mails").Find(nil).All(&mailsCollection.Data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(mailsCollection)
}

func (uc MailController) RetrieveMail(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	mr := models.MailResource{}

	id := params.ByName("id")
	err := uc.session.DB("mailer").C("mails").FindId(bson.ObjectIdHex(id)).One(&mr.Data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(mr)
}

func (uc MailController) CreateMail(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*models.MailResource)

	err := uc.session.DB("mailer").C("mails").Insert(body.Data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(body)
}

func (uc MailController) UpdateMail(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	body := context.Get(r, "body").(*models.MailResource)
	body.Data.Id = bson.ObjectIdHex(params.ByName("id"))

	err := uc.session.DB("mailer").C("mails").UpdateId(body.Data.Id, body.Data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}

func (uc MailController) DeleteMail(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	id := params.ByName("id")

	err := uc.session.DB("mailer").C("mails").RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(204)
	w.Write([]byte("\n"))
}
