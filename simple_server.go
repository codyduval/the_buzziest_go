package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"os"
)

type Person struct {
	Name  string
	Email string
}

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
func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, rootForm)
}

const rootForm = `
  <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <title>Your details</title>
        <link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.4.2/pure-min.css">
      </head>
      <body style="margin: 20px;">
        <h2>A Fun Go App on Heroku to access MongoDB on MongoHQ</h2>
        <p>This simple app will fetch the email id of a person, if it's already there in the MongoDB database.</p>
        <p>Please enter a name (example: Stefan Klaste)</p>
        <form action="/display" method="post" accept-charset="utf-8" class="pure-form">
          <input type="text" name="name" placeholder="name" />
          <input type="submit" value=".. and query database!" class="pure-button pure-button-primary"/>
    </form>
      </body>
    </html>
`

var displayTemplate = template.Must(template.New("display").Parse(displayTemplateHTML))

func display(w http.ResponseWriter, r *http.Request) {
	// In the open command window set the following for Heroku:
	// heroku config:set MONGOHQ_URL=
	// mongodb://IndianGuru:password@troup.mongohq.com:10080/godata
	uri := os.Getenv("MONGOHQ_URL")
	if uri == "" {
		fmt.Println("no connection string provided, using local db")
		uri = "localhost:27017/test"
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB("test").C("users")

	result := Person{}

	collection.Find(bson.M{"name": r.FormValue("name")}).One(&result)

	if result.Email != "" {
		errn := displayTemplate.Execute(w, "The email id you wanted is: "+result.Email)
		if errn != nil {
			http.Error(w, errn.Error(), http.StatusInternalServerError)
		}
	} else {
		displayTemplate.Execute(w, "Sorry... The email id you wanted does not exist.")
	}
}

const displayTemplateHTML = ` 
<!DOCTYPE html>
  <html>
    <head>
      <meta charset="utf-8">
      <title>Results</title>
      <link rel="stylesheet" href="http://yui.yahooapis.com/pure/0.4.2/pure-min.css">
    </head>
    <body>
      <h2>A Fun Go App on Heroku to access MongoDB on MongoHQ</h2>
      <p><b>{{html .}}</b></p>
      <p><a href="/">Start again!</a></p>
    </body>
  </html>
`

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/display", display)
	fmt.Println("listening...")
	err := http.ListenAndServe(GetPort(), nil)
	if err != nil {
		panic(err)
	}
}
