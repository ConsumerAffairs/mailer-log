package main

// http://jsonapi.org
// TODO: Schema valiadation ?
// TODO: Pagination - http://jsonapi.org/format/#fetching-pagination
// TODO: Filtering - http://jsonapi.org/recommendations/#filtering
// TODO: Implement fields type, links and etc.- http://jsonapi.org/format/#document-structure
// TODO: Config file
// TODO: Disable LoggingHandler if nto DEBUG

import (
	"github.com/ConsumerAffairs/mailer_microservice/handlers"
	"github.com/ConsumerAffairs/mailer_microservice/middleware"
	"github.com/ConsumerAffairs/mailer_microservice/models"
	"github.com/ConsumerAffairs/mailer_microservice/router"
	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"gopkg.in/mgo.v2"
	"net/http"
)

var (
	session    *mgo.Session
	collection *mgo.Collection
)

func FileServerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {
	mc := handlers.NewMailController(getSession())
	middleware := middleware.Middleware{}
	commonHandlers := alice.New(context.ClearHandler,
		//middleware.LoggingHandler,
		middleware.RecoverHandler,
		middleware.AcceptHandler)

	router := router.NewRouter()
	router.Get("/mails", commonHandlers.ThenFunc(mc.ListMail))
	router.Post("/mails", commonHandlers.Append(middleware.ContentTypeHandler, middleware.BodyHandler(models.MailResource{})).ThenFunc(mc.CreateMail))

	router.Put("/mail/:id", commonHandlers.Append(middleware.ContentTypeHandler, middleware.BodyHandler(models.MailResource{})).ThenFunc(mc.UpdateMail))
	router.Get("/mail/:id", commonHandlers.ThenFunc(mc.RetrieveMail))
	router.Delete("/mail/:id", commonHandlers.ThenFunc(mc.DeleteMail))
	router.Router.ServeFiles("/static/*filepath", http.Dir("static"))

	http.ListenAndServe(":8080", router)
}

// Get mongodb session
func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	s.SetSafe(&mgo.Safe{})

	// Ensure indexes
	index := mgo.Index{
		Key:        []string{"sent_at"},
		Background: true,
	}

	err = s.DB("mailer").C("mails").EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	// Ensure full text index
	textindex := mgo.Index{
		Key:        []string{"$text:$**"},
		Background: true,
	}

	err = s.DB("mailer").C("mails").EnsureIndex(textindex)
	if err != nil {
		panic(err)
	}

	return s
}
