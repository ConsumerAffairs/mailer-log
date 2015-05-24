package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ConsumerAffairs/mailer-log/models"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	mailsCollection := []models.Mail{}

	u, err := url.Parse(r.URL.String())
	values, err := url.ParseQuery(u.RawQuery)

	page := 0
	if value, ok := values["page"]; ok {
		page, err = strconv.Atoi(value[0])
		if page > 0 {
			page -= 1
		} else if err != nil {
			page = 0
		}
	}

	per_page := 10
	if value, ok := values["per_page"]; ok {
		per_page, err = strconv.Atoi(value[0])
		if err != nil {
			per_page = 10
		}
	}

	skip := page * per_page

	query := uc.session.DB("mailer").C("mails").Find(nil)

	count, err := query.Count()
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", count)

	err = query.Skip(skip).Limit(per_page).All(&mailsCollection)
	if err != nil {
		panic(err)
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%v", count))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mailsCollection)
}

func (uc MailController) RetrieveMail(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	mr := models.Mail{}

	id := params.ByName("id")
	err := uc.session.DB("mailer").C("mails").FindId(bson.ObjectIdHex(id)).One(&mr)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mr)
}

func (uc MailController) CreateMail(w http.ResponseWriter, r *http.Request) {
	body := context.Get(r, "body").(*models.Mail)

	err := uc.session.DB("mailer").C("mails").Insert(body)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(body)
}

func (uc MailController) UpdateMail(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	body := context.Get(r, "body").(*models.Mail)
	body.Id = bson.ObjectIdHex(params.ByName("id"))

	err := uc.session.DB("mailer").C("mails").UpdateId(body.Id, body)
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
