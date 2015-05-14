package main

// http://jsonapi.org
// TODO: Schema valiadation ?
// TODO: Pagination - http://jsonapi.org/format/#fetching-pagination
// TODO: Filtering - http://jsonapi.org/recommendations/#filtering
// TODO: Implement fields type, links and etc.- http://jsonapi.org/format/#document-structure
// TODO: Config file
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
    "stathat.com/c/jconfig"
)

var (
	session    *mgo.Session
	collection *mgo.Collection
)

type Config struct {
    debug       bool
    host        string
    mongodb     string

}




func FileServerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func main() {

    config := jconfig.LoadConfig("config.json")

	mc := handlers.NewMailController(getSession(config.GetString("mongodb")))
	middleware := middleware.Middleware{}
	commonHandlers := alice.New(context.ClearHandler)

    if config.GetBool("debug") {
        commonHandlers = commonHandlers.Append(middleware.LoggingHandler)
    }

    commonHandlers = commonHandlers.Append(middleware.RecoverHandler,
                                            middleware.AcceptHandler)

	router := router.NewRouter()
	router.Get("/mails", commonHandlers.ThenFunc(mc.ListMail))
	router.Post("/mails", commonHandlers.Append(middleware.ContentTypeHandler, middleware.BodyHandler(models.MailResource{})).ThenFunc(mc.CreateMail))

	router.Put("/mail/:id", commonHandlers.Append(middleware.ContentTypeHandler, middleware.BodyHandler(models.MailResource{})).ThenFunc(mc.UpdateMail))
	router.Get("/mail/:id", commonHandlers.ThenFunc(mc.RetrieveMail))
	router.Delete("/mail/:id", commonHandlers.ThenFunc(mc.DeleteMail))
	router.Router.ServeFiles("/static/*filepath", http.Dir("static"))

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
