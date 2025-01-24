package crawler

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Crawler struct {
	colly *colly.Collector
	db    *sql.DB
}

type Brand struct {
	Name        string
	Description string
	Strength    string
	Country     string
	Wrapper     string
	Shapes      string
	Products    []Product
}

type Product struct {
	Name      string
	Packaging []Packaging
}

type Packaging struct {
	Name  string
	MSRP  string
	Price string
}

func NewCrawler(collector *colly.Collector, db *sql.DB) *Crawler {
	collector.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"
	return &Crawler{
		colly: collector,
		db:    db,
	}
}

func (crawler *Crawler) Start() {

	crawler.colly.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %s", r.URL)
	})

	crawler.colly.OnError(func(r *colly.Response, err error) {
		if r.StatusCode != 200 {
			log.Printf("Error occured trying to reach %s : %s", r.Request.URL, err)
			return
		}
		if err != nil {
			log.Printf("Something went wrong visiting %s : %s", r.Request.URL, err)
			return
		}
	})

	crawler.colly.OnHTML("div#top", func(element *colly.HTMLElement) {
		productLink := element.ChildAttr("a", "href")
		crawler.scrapeCigarData(productLink)
		// element.ForEach("li.category_link", func(i int, h *colly.HTMLElement) {
		// 	cigarBrandLink := h.ChildAttr("a", "href")
		// 	if strings.Contains(cigarBrandLink, "https") {

		// 		go func(link string) {
		// 			crawler.scrapeCigarData(link)
		// 		}(cigarBrandLink)
		// 	}
		// })
	})

	crawler.colly.Visit("https://www.holts.com/cigars/all-cigar-brands.html")
}

func (crawler *Crawler) scrapeCigarData(link string) {
	crawler.colly.OnRequest(func(r *colly.Request) {
		log.Printf("Visiting %s", r.URL)
	})

	crawler.colly.OnError(func(r *colly.Response, err error) {
		if r.StatusCode != 200 {
			log.Printf("Error occured trying to reach %s : %s", r.Request.URL, err)
			return
		}
		if err != nil {
			log.Printf("Something went wrong visiting %s : %s", r.Request.URL, err)
			return
		}
	})

	cigarBrand := &Brand{}

	crawler.colly.OnHTML("div.col-main", func(h *colly.HTMLElement) {
		cigarBrand.Name = h.ChildText("h1")
		cigarBrand.Description = h.ChildText("div.product-description p")
		cigarBrand.Strength = h.ChildText("div.strength-o-meter div.value")
		cigarBrand.Country = strings.ReplaceAll(h.ChildText("div.pdp-cigar-details > ul > li:nth-child(2) > span > span"), " ", "")
		cigarBrand.Wrapper = strings.ReplaceAll(h.ChildText("div.pdp-cigar-details > ul > li:last-child > span > span"), " ", "")
		cigarBrand.Shapes = strings.TrimSpace(strings.Split(h.ChildText("div.pdp-cigar-details > div.sizes"), ":")[1])

		product := Product{}
		h.ForEach("tr", func(i int, h *colly.HTMLElement) {

			if h.ChildText("div.name-wrapper > div.name") != "" {
				product.Name = h.ChildText("div.name-wrapper > div.name")
				if h.ChildText("td.tpacking") != "" {
					product.Packaging = append(product.Packaging, Packaging{Name: h.ChildText("td.tpacking"), MSRP: h.ChildText("td.tmsrp"), Price: h.ChildText("td.tprice")})
				}
			} else {
				if h.ChildText("td.tpacking") != "" {
					product.Packaging = append(product.Packaging, Packaging{Name: h.ChildText("td.tpacking"), MSRP: h.ChildText("td.tmsrp"), Price: h.ChildText("td.tprice")})
				}
			}

			if h.Attr("class") == "last-row" {
				cigarBrand.Products = append(cigarBrand.Products, product)
				product = Product{}
			}

		})
	})

	crawler.colly.OnScraped(func(r *colly.Response) {
		_, err := crawler.db.Query(
			fmt.Sprintf(`
			INSERT INTO brands (name,description,strength,country,wrapper,shapes) 
			VALUES('%s','%s','%s','%s','%s','%s');`, cigarBrand.Name, cigarBrand.Description, cigarBrand.Strength, cigarBrand.Country, cigarBrand.Wrapper, cigarBrand.Shapes))
		if err != nil {
			log.Printf("Error occured inserting cigar data into table : %s", err)
		}
		log.Println(cigarBrand.Name)
		log.Println(cigarBrand.Description)
		log.Println(cigarBrand.Strength)
		log.Println(cigarBrand.Wrapper)
		log.Println(cigarBrand.Country)
		log.Println(cigarBrand.Shapes)

		for _, prod := range cigarBrand.Products {
			log.Println(prod)
		}

	})
	crawler.colly.Visit(link)
}
