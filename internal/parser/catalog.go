package parser

import (
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

type Category struct {
	Slug       string
	URL        string
	LeicaValue string
}

var Categories = []Category{
	{
		Slug:       "3d-scanners",
		URL:        "https://www.geooptic.ru/catalog/3d-skanirovanie",
		LeicaValue: "2",
	},
	{
		Slug:       "total-stations",
		URL:        "https://www.geooptic.ru/catalog/taheometry",
		LeicaValue: "0",
	},
	{
		Slug:       "gnss-receivers",
		URL:        "https://www.geooptic.ru/catalog/gnss-priemniki",
		LeicaValue: "2",
	},
	{
		Slug:       "laser-levels",
		URL:        "https://www.geooptic.ru/catalog/lazernye-urovni",
		LeicaValue: "3",
	},
	{
		Slug:       "levels",
		URL:        "https://www.geooptic.ru/catalog/niveliry",
		LeicaValue: "1",
	},
	{
		Slug:       "laser-rangefinders",
		URL:        "https://www.geooptic.ru/catalog/dalnomery",
		LeicaValue: "0",
	},
	{
		Slug:       "ndt-instruments",
		URL:        "https://www.geooptic.ru/catalog/pribory-nerazrushayushchego-kontrolya",
		LeicaValue: "11",
	},
	{
		Slug:       "software",
		URL:        "https://www.geooptic.ru/catalog/soft",
		LeicaValue: "0",
	},
	{
		Slug:       "accessories",
		URL:        "https://www.geooptic.ru/catalog/aksessuary",
		LeicaValue: "1",
	},
}

func (p *Parser) GetProductLinks(category Category) ([]string, error) {
	var links []string

	js := `
		Array.from(document.querySelectorAll(".productList .productItem > a.img"))
			.map(a => a.href);
	`

	leicaSelector := fmt.Sprintf(
		`input[type="checkbox"][value="%s"]`,
		category.LeicaValue,
	)

	err := chromedp.Run(
		p.ctx,

		chromedp.Navigate(category.URL),

		chromedp.WaitVisible(leicaSelector),
		chromedp.Sleep(2*time.Second),

		chromedp.Evaluate(
			fmt.Sprintf(`
				const el = document.querySelector('input[type="checkbox"][value="%s"]');
				el.checked = true;
				el.dispatchEvent(new Event('change', { bubbles: true }));
			`, category.LeicaValue),
			nil,
		),

		chromedp.Sleep(2*time.Second),

		// Раскрываем все товары категории
		chromedp.Evaluate(`
			const btn = document.querySelector(
				'span[title="Не запариваться и показать всё"]'
			);
			if (btn) btn.click();
		`, nil),

		chromedp.Sleep(2*time.Second),

		chromedp.Evaluate(js, &links),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"не удалось получить ссылки для категории %s: %w",
			category.Slug,
			err,
		)
	}

	unique := make(map[string]bool)
	var result []string

	for _, link := range links {
		if unique[link] {
			continue
		}

		unique[link] = true
		result = append(result, link)
	}

	return result, nil
}
