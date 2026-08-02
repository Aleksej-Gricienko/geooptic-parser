package main

import (
	"log"

	"geooptic-parser-chromedp/internal"
)

func main() {
	ctx, cancel := internal.NewBrowser()
	defer cancel()

	parser := internal.NewParser(ctx)

	if err := parser.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Парсинг завершён.")
}
