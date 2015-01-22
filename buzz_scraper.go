package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/xmlpath.v2"
	"log"
	"net/http"
	"time"
)

type (
	BlogPost struct {
		ID      bson.ObjectId `bson:"_id,omitempty"`
		Date    time.Time     `bson:"date"`
		Uri     string        `bson:"uri"`
		Content string        `bson:"content"`
		Title   string        `bson:"title"`
		Source  string        `bson:"source"`
		Guid    string        `bson:"guid"`
	}
)

func main() {

	url := "http://blog.codyduval.com"
	// Xpath to the title of the first blog post
	node := "/html/body/div/div/div/div/div[1]/div[1]/a/h3"

	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
	}

	path, err := xmlpath.Compile(node)
	if err != nil {
		fmt.Println("Didn't work!")
		log.Fatal(err)
	}

	root, err := xmlpath.ParseHTML(r.Body)
	if err != nil {
		fmt.Println("Bummer, I couldn't parse that blob of html!")
		log.Fatal(err)
	}

	if value, ok := path.String(root); ok {
		fmt.Println("Found:", value)
	} else {
		fmt.Println("Sorry, couldn't find anything at the xpath node: " + node)
	}

}
