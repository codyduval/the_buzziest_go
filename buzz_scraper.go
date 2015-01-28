package main

import (
	"fmt"
	"gopkg.in/xmlpath.v2"
	"log"
	"net/http"
	"time"
)

func main() {
	var i int

	fmt.Println("How many pages to grab?")
	_, err := fmt.Scanf("%d", &i)
	if err != nil {
		fmt.Println(err)
	}
	numOfPages := i
	fmt.Println("Concurrent(1) or Sequential(2):")
	_, err = fmt.Scanf("%d", &i)
	if err != nil {
		fmt.Println(err)
	}
	conOrSeq := i

	if conOrSeq == 1 {

		start := time.Now()

		responses := fetchPages(numOfPages)
		for i := 1; i <= numOfPages; i++ {
			response := <-responses
			for node := 2; node <= 6; node++ {
				xpath := fmt.Sprintf("//*[@id='dispatch']/div[1]/div[6]/div[4]/div[%v]/div/h1", node)

				name := scrapePage(response, xpath)
				fmt.Println("Found:", name)
			}
		}
		elapsed := time.Since(start)
		log.Printf("Concurrent took %s", elapsed)
	} else {
		start := time.Now()
		for page := 1; page <= numOfPages; page++ {
			url := fmt.Sprintf("http://www.tastingtable.com/dispatch/national/dispatch?applyFilter=true&filterEditionId=1&filterNeighborhoodId=0&filterDropDown1=opened&requestedPage=%v", page)
			response := fetchPage(url)
			root := parsePage(response)
			for node := 2; node <= 6; node++ {
				xpath := fmt.Sprintf("//*[@id='dispatch']/div[1]/div[6]/div[4]/div[%v]/div/h1", node)

				name := scrapePage(root, xpath)
				fmt.Println("Found:", name)
			}

		}
		elapsed := time.Since(start)
		log.Printf("Sequential took %s", elapsed)
	}

}

func fetchPages(numOfPages int) <-chan *xmlpath.Node {
	ch := make(chan *xmlpath.Node) // buffered
	for page := 1; page <= numOfPages; page++ {
		url := fmt.Sprintf("http://www.tastingtable.com/dispatch/national/dispatch?applyFilter=true&filterEditionId=1&filterNeighborhoodId=0&filterDropDown1=opened&requestedPage=%v", page)
		go func(url string) {
			fmt.Printf("Fetching %s \n", url)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Something went wrong!")
				log.Fatal(err)
			}
			root := parsePage(resp)
			ch <- root
		}(url)
	}
	return ch
}

func fetchPage(url string) *http.Response {
	fmt.Printf("Fetching %s \n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Something went wrong!")
		log.Fatal(err)
	}
	return resp
}

func parsePage(r *http.Response) *xmlpath.Node {

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
