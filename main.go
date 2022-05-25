package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

func main() {

	log.Print("starting server...")
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	var prices []Price

	c := colly.NewCollector()

	for i := 0; i < 10; i++ {

		c.OnHTML("div.tabstable tbody", func(e *colly.HTMLElement) {
			e.ForEach("tr", func(j int, el *colly.HTMLElement) {
				if j == 0 || j == 1 {
					return
				}

				price := Price{}
				if el.ChildText("td:nth-child(1)") != "" {
					price.Sno = el.ChildText("td:nth-child(1)")
				}
				if el.ChildText("td:nth-child(2)") != "" {
					price.Date = el.ChildText("td:nth-child(2)")
				}
				if el.ChildText("td:nth-child(3)") != "" {
					price.Market = el.ChildText("td:nth-child(3)")
				}
				if el.ChildText("td:nth-child(4)") != "" {
					price.Type = el.ChildText("td:nth-child(4)")
				}
				if el.ChildText("td:nth-child(5)") != "" {
					price.Price = el.ChildText("td:nth-child(5)")
				}

				prices = append(prices, price)

			})
		})
		url := fmt.Sprintf("http://www.indianspices.com/indianspices/marketing/price/domestic/daily-price-large.html?page=%s", strconv.Itoa(i))
		c.Visit(url)
	}
	cardamom := Cardamom{
		Prices: &prices,
	}
	b, err := json.Marshal(cardamom)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprint(w, string(b))

}

type Cardamom struct {
	Prices *[]Price `json:"prices,omitempty"`
}

type Price struct {
	Sno    string `json:"sno,omitempty"`
	Date   string `json:"date,omitempty"`
	Market string `json:"market,omitempty"`
	Type   string `json:"type,omitempty"`
	Price  string `json:"price,omitempty"`
}
