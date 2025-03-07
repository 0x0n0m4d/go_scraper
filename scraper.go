package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Product struct {
	Url, Name, Price, Image string
}

func main() {
	collector := colly.NewCollector(
		colly.AllowedDomains("www.aliexpress.com"),
	)

	collector.UserAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0"

	var products []Product

	collector.OnError(func(_ *colly.Response, err error) {
		log.Fatalln(err)
	})

	collector.OnHTML("div.search-item-card-wrapper-gallery", func(e *colly.HTMLElement) {
		product := Product{}

		product.Url = e.ChildAttr("a.search-card-item", "href")
		product.Name = e.ChildAttr(".lq_ae", "title")
		product.Price = e.ChildText("div.lq_j3 > span")
		product.Image = e.ChildAttr("img.l9_be", "src")

		products = append(products, product)
	})

	collector.OnScraped(func(r *colly.Response) {
		file, err := os.Create("products.csv")
		if err != nil {
			log.Fatalln("Failed to create output CSV file", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)

		headers := []string{
			"Url",
			"Name",
			"Price",
			"Image",
		}
		writer.Write(headers)

		fmt.Printf("Scraped %d products!\n", len(products))
		for _, product := range products {
			record := []string{
				"https:" + product.Url,
				product.Name,
				product.Price,
				"https:" + product.Image,
			}

			writer.Write(record)
		}
		defer writer.Flush()
	})

	collector.Visit("https://www.aliexpress.com/w/wholesale-cellphone.html")
}
