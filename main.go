package main

import (
	crawler "cigar-scraper/Crawler"
	"cigar-scraper/config"
	"database/sql"
	"fmt"
	"log"

	"github.com/gocolly/colly"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", cfg.DB.URL, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Printf("database ping failed : %s", err)
	}
	_, err = db.Query(`CREATE TABLE Brands (
			Name VARCHAR(50) PRIMARY KEY,
			Description TEXT,
			Strength VARCHAR(20),
			Country VARCHAR(50),
			Wrapper VARCHAR(100),
			Shapes VARCHAR(100),
			Products text ARRAY
	)`)
	if err != nil {
		log.Printf("Error creating new table : %s", err)
	}
	log.Println("Table successfully created")
	c := colly.NewCollector()
	crawler := crawler.NewCrawler(c, db)
	crawler.Start()
}
