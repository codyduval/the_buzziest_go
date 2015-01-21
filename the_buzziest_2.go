package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type (
	BuzzSource struct {
		ID          bson.ObjectId `bson:"_id,omitempty"`
		Name        string        `bson:"name"`
		Uri         string        `bson:"uri"`
		BuzzWeight  float64       `bson:"buzz_weight"`
		XPathNodes  string        `bson:"x_path_nodes"`
		SourceType  string        `bson:"source_type"`
		DecayFactor float64       `bson:"decay_factor"`
		City        string        `bson:"city"`
		Timestamp   time.Time     `bson:"time_stamp"`
	}

	BuzzPost struct {
		ID      bson.ObjectId `bson:"_id,omitempty"`
		Date    time.Time     `bson:"date"`
		Uri     string        `bson:"uri"`
		Content string        `bson:"content"`
		Title   string        `bson:"title"`
		Source  BuzzSource    `bson:"source"`
		Guid    string        `bson:"guid"`
	}
)

const buzzSourceForm = `
  <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <title>Your details</title>
        <link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.4.2/pure-min.css">
      </head>
      <body style="margin: 20px;">
        <h2>Enter a Buzz Source</h2>
        <form action="/buzz_sources/create" method="post" accept-charset="utf-8" class="pure-form">
        	<p>Name</p>
          <input type="text" name="name" placeholder="name" />
        	<p>URI</p>
          <input type="text" name="uri" placeholder="uri" />
        	<p>Weight</p>
          <input type="text" name="buzz_weight" placeholder="weight" />
        	<p>XPath Node</p>
          <input type="text" name="x_path_nodes" placeholder="x path nodes" />
        	<p>Type</p>
          <input type="text" name="source_type" placeholder="blog, twitter, etc" />
        	<p>Decay Factor</p>
          <input type="text" name="decay_factor" placeholder="default = 0.906" />
        	<p>City</p>
          <input type="text" name="city" placeholder="nyc, la, or sf" />
          <input type="submit" value="... submit!" class="pure-button pure-button-primary"/>
    </form>
      </body>
    </html>
`

func main() {
	// We need this object to establish a session to our MongoDB.
	rend := render.New()
	dbSession := getSession()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HomeHandler)

	// BuzzSources collection
	buzzSources := router.PathPrefix("/buzz_sources").Subrouter()
	buzzSources.Methods("GET").Path("/new").HandlerFunc(BuzzSourcesNewHandler)
	buzzSources.Methods("GET").HandlerFunc(BuzzSourcesIndexHandler(rend))
	buzzSources.Methods("POST").Path("/create").HandlerFunc(BuzzSourcesCreateHandler(dbSession))

	// BuzzSource (singular)
	buzzSource := router.PathPrefix("/buzz_source/{id}").Subrouter()
	buzzSource.Methods("GET").Path("/edit").HandlerFunc(BuzzSourceEditHandler)
	buzzSource.Methods("GET").HandlerFunc(BuzzSourceShowHandler)
	buzzSource.Methods("PUT", "POST").HandlerFunc(BuzzSourceUpdateHandler)
	buzzSource.Methods("DELETE").HandlerFunc(BuzzSourceDeleteHandler)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}

var (
	mgoSession *mgo.Session
)

func NewBuzzSource() *BuzzSource {
	return &BuzzSource{
		Timestamp: time.Now(),
	}
}

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial("localhost")
		if err != nil {
			panic(err) // no, not really
		}
	}
	return mgoSession.Copy()
}

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Home")
}

func BuzzSourcesIndexHandler(rend *render.Render) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rend.JSON(rw, http.StatusOK, map[string]string{"hello": "json"})
	}
}

func BuzzSourcesCreateHandler(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {

		collection := session.DB("test").C("buzz_sources")

		entry := NewBuzzSource()
		entry.Name = r.FormValue("name")
		fmt.Println(entry)

		if err := collection.Insert(entry); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(rw, r, "/buzz_sources/new", http.StatusTemporaryRedirect)

		session.Close()
	}
}

func BuzzSourcesNewHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, buzzSourceForm)
}

func BuzzSourceShowHandler(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintln(rw, "showing buzz_source ", id)
}

func BuzzSourceUpdateHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "buzz_source update")
}

func BuzzSourceDeleteHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "buzz_source delete")
}

func BuzzSourceEditHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "buzz_source edit")
}
