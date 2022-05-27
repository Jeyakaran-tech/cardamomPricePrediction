package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	logger "log"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine/log"

	"github.com/gocolly/colly"
	"google.golang.org/appengine"
)

func main() {

	logger.Print("starting server...")
	http.HandleFunc("/cardamom", handler)
	appengine.Main()

}

func handler(w http.ResponseWriter, r *http.Request) {

	bucket := "development-cardamomprice"
	fileName := "cardamom-jk-go"
	// buf := &bytes.Buffer{}
	status := Status{
		Code:        "8200",
		Description: "Success",
	}

	var prices []Price
	c := colly.NewCollector()

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to create client: %v", err)
		return
	}
	defer client.Close()

	wc := client.Bucket(bucket).Object(fileName).NewWriter(ctx)
	//PROGRAMMING_LOGIC_FOR_DATA_EXTRACTION
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
	byteResponse, err := json.Marshal(cardamom)
	if err != nil {
		fmt.Println(err)
		return
	}
	statusResponse, errStatus := json.Marshal(status)
	if errStatus != nil {
		fmt.Println(errStatus)
		return
	}

	w.WriteHeader(http.StatusOK)
	//PROGRAMMING_LOGIC_FINISHED
	wc.ContentType = "application/json"

	io.Copy(wc, bytes.NewReader(byteResponse))

	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := wc.Write([]byte(statusResponse)); err != nil {
		log.Errorf(ctx, "createFile: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
		return
	}
	if err := wc.Close(); err != nil {
		log.Errorf(ctx, "createFile: unable to close bucket %q, file %q: %v", bucket, fileName, err)
		return
	}

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

type Status struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}
