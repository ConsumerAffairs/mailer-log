package main

// http://jsonapi.org
// TODO: Schema valiadation ?
// TODO: Pagination - http://jsonapi.org/format/#fetching-pagination
// TODO: Filtering - http://jsonapi.org/recommendations/#filtering
// TODO: Implement fields type, links and etc.- http://jsonapi.org/format/#document-structure
// TODO: Disable LoggingHandler if nto DEBUG

import (
	"github.com/ConsumerAffairs/mailer-log/handlers"
	"github.com/ConsumerAffairs/mailer-log/middleware"
	"github.com/ConsumerAffairs/mailer-log/models"
	"github.com/ConsumerAffairs/mailer-log/router"
	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"stathat.com/c/jconfig"
)

var (
	session    *mgo.Session
	collection *mgo.Collection
)

func main() {

	var config *jconfig.Config
	if _, err := os.Stat("config.local.json"); os.IsNotExist(err) {
		config = jconfig.LoadConfig("config.json")
	} else {
		config = jconfig.LoadConfig("config.local.json")
	}

	mc := handlers.NewMailController(getSession(config.GetString("mongodb")))
	middleware := middleware.Middleware{}
	commonHandlers := alice.New(context.ClearHandler)

	if config.GetBool("debug") {
		commonHandlers = commonHandlers.Append(middleware.LoggingHandler)
	}

	commonHandlers = commonHandlers.Append(
		middleware.RecoverHandler,
		middleware.AcceptHandler)

	router := router.NewRouter()
	router.Get("/mails", commonHandlers.ThenFunc(mc.ListMail))
	router.Post("/mails", commonHandlers.Append(middleware.ContentTypeHandler,
		middleware.BodyHandler(models.Mail{})).ThenFunc(mc.CreateMail))

	router.Put("/mails/:id", commonHandlers.Append(middleware.ContentTypeHandler,
		middleware.BodyHandler(models.Mail{})).ThenFunc(mc.UpdateMail))
	router.Get("/mails/:id", commonHandlers.ThenFunc(mc.RetrieveMail))
	router.Delete("/mails/:id", commonHandlers.ThenFunc(mc.DeleteMail))
	router.Router.NotFound = http.FileServer(http.Dir(config.GetString("servePath"))).ServeHTTP

	http.ListenAndServe(config.GetString("host"), router)
}

// Get mongodb session
func getSession(url string) *mgo.Session {
	s, err := mgo.Dial("")
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
