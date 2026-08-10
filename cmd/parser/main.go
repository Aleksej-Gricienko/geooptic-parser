package main

import (
	"log"

	"geooptic-parser-chromedp/internal/parser"
)

func main() {
	ctx, cancel := parser.NewBrowser()
	defer cancel()

	p := parser.NewParser(ctx)

	if err := p.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Парсинг завершён.")
}
