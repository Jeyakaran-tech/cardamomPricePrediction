package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	logger "log"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine/log"

	"github.com/gocolly/colly"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
)

func main() {

	logger.Print("starting server...")
	http.HandleFunc("/cardamom", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	logger.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal(err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)

	//[START get_default_bucket]
	// Use `dev_appserver.py --default_gcs_bucket_name GCS_BUCKET_NAME`
	// when running locally.
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get default GCS bucket name: %v", err)
	}
	//[END get_default_bucket]

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to create client: %v", err)
		return
	}
	defer client.Close()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Demo GCS Application running from Version: %v\n", appengine.VersionID(ctx))
	fmt.Fprintf(w, "Using bucket name: %v\n\n", bucket)

	buf := &bytes.Buffer{}
	d := &demo{
		w:          buf,
		ctx:        ctx,
		client:     client,
		bucket:     client.Bucket(bucket),
		bucketName: bucket,
	}

	n := "cardamom-jk-go"
	d.createFile(n)

	if d.failed {
		w.WriteHeader(http.StatusInternalServerError)
		buf.WriteTo(w)
		fmt.Fprintf(w, "\nfailed.\n")
	} else {
		w.WriteHeader(http.StatusOK)
		buf.WriteTo(w)
		fmt.Fprintf(w, "\nsucceeded.\n")
	}

}

func (d *demo) createFile(fileName string) {
	var prices []Price
	m := make(map[string]string)
	c := colly.NewCollector()
	fmt.Fprintf(d.w, "Creating file /%v/%v\n", d.bucketName, fileName)

	wc := d.bucket.Object(fileName).NewWriter(d.ctx)
	wc.ContentType = "application/json"

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
	err1 := json.Unmarshal(b, &m)
	if err1 != nil {
		logger.Fatal()
	}
	wc.Metadata = m

	if _, err := wc.Write([]byte("abcde\n")); err != nil {
		d.errorf("createFile: unable to write data to bucket %q, file %q: %v", d.bucketName, fileName, err)
		return
	}
	if _, err := wc.Write([]byte(strings.Repeat("f", 1024*4) + "\n")); err != nil {
		d.errorf("createFile: unable to write data to bucket %q, file %q: %v", d.bucketName, fileName, err)
		return
	}
	if err := wc.Close(); err != nil {
		d.errorf("createFile: unable to close bucket %q, file %q: %v", d.bucketName, fileName, err)
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

type demo struct {
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle

	w   io.Writer
	ctx context.Context
	// failed indicates that one or more of the demo steps failed.
	failed bool
}

func (d *demo) errorf(format string, args ...interface{}) {
	d.failed = true
	fmt.Fprintln(d.w, fmt.Sprintf(format, args...))
	log.Errorf(d.ctx, format, args...)
}
