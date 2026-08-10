package main

import (
	"log"

	"geooptic-parser-chromedp/internal/importer"
)

func main() {
	db, err := importer.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.DB.Close()

	err = db.ImportCategoryFile("data/levels.json")

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

}
