package parser

import (
	"context"
	"log"

	"geooptic-parser-chromedp/models"
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

	for _, category := range Categories {

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

		// ↓ Скачиваем изображения
		// if err := DownloadImages(products); err != nil {
		//     return err
		// }

		// ↓ Сохраняем документы
		if err := DownloadDocuments(products, category.Slug); err != nil {
			return err
		}

		// ↓ Сохраняем товары
		if err := SaveProducts(products, category.Slug); err != nil {
			return err
		}
	}

	log.Println("Все товары успешно обработаны.")

	return nil
}
