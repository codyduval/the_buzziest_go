package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/buzz_sources", BuzzSourcesHandler)
	http.Handle("/", r)
}
