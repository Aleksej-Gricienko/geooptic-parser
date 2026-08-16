package parser

import (
	"context"

	"geooptic-parser-chromedp/models"
	"log"
)

type Parser struct {
	ctx context.Context
}

func NewParser(ctx context.Context) *Parser {
	return &Parser{
		ctx: ctx,
	}
}
func (p *Parser) Run() error {

	for _, category := range []Category{
		{
			Slug:         "accessories",
			URL:          "https://www.geooptic.ru/catalog/aksessuary",
			FilterValue:  "1",
			Manufacturer: "Leica Geosystems",
		},
	} {

		links, err := p.GetProductLinks(category)
		if err != nil {
			return err
		}

		log.Printf("Найдено ссылок: %d\n", len(links))

		var products []models.Product

		for _, link := range links {
			product, err := p.ParseProduct(link)
			if err != nil {
				log.Println(err)
				continue
			}

			products = append(products, product)
		}

		// Скачиваем изображения
		if err := DownloadImages(products, category.Slug); err != nil {
			return err
		}

		// Скачиваем документы
		if err := DownloadDocuments(products, category.Slug); err != nil {
			return err
		}

		// Сохраняем JSON
		if err := SaveProducts(products, category.Slug); err != nil {
			return err
		}
	}

	log.Println("Все товары успешно обработаны.")

	return nil
}
