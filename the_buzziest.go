package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	MongoDBHosts = "ds035428.mongolab.com:35428"
	AuthDatabase = "goinggo"
	AuthUserName = "guest"
	AuthPassword = "welcome"
	TestDatabase = "goinggo"
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

// Get the Port from the environment so we can run on Heroku
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "9090"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
func newBuzzSource(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, buzzSourceForm)
}

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
        <form action="/display" method="post" accept-charset="utf-8" class="pure-form">
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
