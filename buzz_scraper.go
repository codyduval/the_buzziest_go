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

	var url string
	var xpath string
	pages_to_pull := 5

	for page := 1; page <= pages_to_pull; page++ {
		url = fmt.Sprintf("http://www.tastingtable.com/dispatch/national/dispatch?applyFilter=true&filterEditionId=1&filterNeighborhoodId=0&filterDropDown1=opened&requestedPage=%v", page)

		response := make(chan *http.Response)
		go fetchPage(url, response)
		root := parsePage(response)

		for node := 2; node <= 6; node++ {
			xpath = fmt.Sprintf("//*[@id='dispatch']/div[1]/div[6]/div[4]/div[%v]/div/h1", node)

			name := scrapePage(root, xpath)
			fmt.Println("Found:", name)

		}
	}

}

func fetchPage(url string, response chan *http.Response) {
	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
	}
	response <- r
}

func parsePage(r chan *http.Response) *xmlpath.Node {

	root, err := xmlpath.ParseHTML(r.Body)
	if err != nil {
		fmt.Println("Bummer, I couldn't parse that blob of html!")
		log.Fatal(err)
	}
	return root

}

func scrapePage(root *xmlpath.Node, xpath string) string {
	path, err := xmlpath.Compile(xpath)
	if err != nil {
		fmt.Println("Didn't work!")
		log.Fatal(err)
	}

	if name, ok := path.String(root); ok {
		return name
	} else {
		return "Sorry, couldn't find anything at the xpath node"
	}

}
