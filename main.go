package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	logger "log"
	"net/http"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gocolly/colly"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func main() {
	logger.Print("starting server...")
	http.HandleFunc("/cardamom", handler)
	appengine.Main()
}
func handler(w http.ResponseWriter, r *http.Request) {

	bucket := "development-cardamomprice"
	fileName := "cardamom-jk-go"
	var emptyBytes []byte
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
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		price := Price{}
		_, err := strconv.ParseInt(e.ChildText("td:nth-child(1)"), 10, 64)
		if err == nil {
			price.Sno = e.ChildText("td:nth-child(1)")
			price.Date = e.ChildText("td:nth-child(2)")
			price.Market = e.ChildText("td:nth-child(3)")
			price.Type = e.ChildText("td:nth-child(4)")
			price.Price = e.ChildText("td:nth-child(5)")
			prices = append(prices, price)
		}
	})
	for i := 1; i <= 10; i++ {
		url := fmt.Sprintf("http://www.indianspices.com/indianspices/marketing/price/domestic/daily-price-large.html?page=%s", strconv.Itoa(i))
		c.Visit(url)
	}
	cardamom := Cardamom{
		Prices: &prices,
		Status: Status{
			Code:        "8200",
			Description: "Success",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	byteResponse, err := json.Marshal(cardamom)
	if err != nil {
		fmt.Println(err)
		return
	}
	// statusResponse, errStatus := json.Marshal(cardamom.Status)
	// if errStatus != nil {
	// 	fmt.Println(errStatus)
	// 	return
	// }

	//PROGRAMMING_LOGIC_FINISHED
	wc.ContentType = "application/json"
	io.Copy(wc, bytes.NewReader(byteResponse))

	if _, err := wc.Write([]byte(emptyBytes)); err != nil {
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
	Status Status   `json:"status,omitempty"`
}
type Status struct {
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

type Price struct {
	Sno    string `json:"sno,omitempty"`
	Date   string `json:"date,omitempty"`
	Market string `json:"market,omitempty"`
	Type   string `json:"type,omitempty"`
	Price  string `json:"price,omitempty"`
}
