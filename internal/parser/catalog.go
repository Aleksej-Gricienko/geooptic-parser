package parser

import (
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

type Category struct {
	Slug         string
	URL          string
	FilterValue  string
	Manufacturer string
}

var Categories = []Category{
	{
		Slug:         "3d-scanners",
		URL:          "https://www.geooptic.ru/catalog/3d-skanirovanie",
		FilterValue:  "2",
		Manufacturer: "Leica Geosystems",
	},
	{
		Slug:         "total-stations",
		URL:          "https://www.geooptic.ru/catalog/taheometry",
		FilterValue:  "0",
		Manufacturer: "Leica Geosystems",
	},
	{
		Slug:         "gnss-receivers",
		URL:          "https://www.geooptic.ru/catalog/gnss-priemniki",
		FilterValue:  "2",
		Manufacturer: "Leica Geosystems",
	},
	{
		Slug:         "laser-levels",
		URL:          "https://www.geooptic.ru/catalog/lazernye-urovni",
		FilterValue:  "3",
		Manufacturer: "Leica Geosystems",
	},
	{
		Slug:         "levels",
		URL:          "https://www.geooptic.ru/catalog/niveliry",
		FilterValue:  "1",
		Manufacturer: "Leica Geosystems",
	},
	{
		Slug:         "laser-rangefinders",
		URL:          "https://www.geooptic.ru/catalog/dalnomery",
		FilterValue:  "0",
		Manufacturer: "Leica Geosystems",
	},
	{
		Slug:         "ndt-instruments",
		URL:          "https://www.geooptic.ru/catalog/pribory-nerazrushayushchego-kontrolya",
		FilterValue:  "11",
		Manufacturer: "Leica Geosystems",
	},
	{
		Slug:         "accessories",
		URL:          "https://www.geooptic.ru/catalog/aksessuary",
		FilterValue:  "1",
		Manufacturer: "Leica Geosystems",
	},
	{
		Slug:         "software",
		URL:          "https://www.geooptic.ru/catalog/soft",
		FilterValue:  "",
		Manufacturer: "",
	},
	{
		Slug:         "gnss-receivers-geomax",
		URL:          "https://www.geooptic.ru/catalog/gnss-priemniki",
		FilterValue:  "13",
		Manufacturer: "GeoMax",
	},
	{
		Slug:         "levels-geomax",
		URL:          "https://www.geooptic.ru/catalog/niveliry",
		FilterValue:  "21",
		Manufacturer: "GeoMax",
	},
	{
		Slug:         "theodolites-geomax",
		URL:          "https://www.geooptic.ru/catalog/teodolity",
		FilterValue:  "5",
		Manufacturer: "GeoMax",
	},
	{
		Slug:         "ndt-instruments-geomax",
		URL:          "https://www.geooptic.ru/catalog/pribory-nerazrushayushchego-kontrolya",
		FilterValue:  "12",
		Manufacturer: "GeoMax",
	},
	{
		Slug:         "accessories-geomax",
		URL:          "https://www.geooptic.ru/catalog/aksessuary",
		FilterValue:  "16",
		Manufacturer: "GeoMax",
	},
	{
		Slug:         "total-stations-geomax",
		URL:          "https://www.geooptic.ru/catalog/taheometry",
		FilterValue:  "2",
		Manufacturer: "GeoMax",
	},
}

func (p *Parser) GetProductLinks(category Category) ([]string, error) {
	var links []string

	js := `
		Array.from(document.querySelectorAll(".productList .productItem > a.img"))
			.map(a => a.href);
	`

	actions := []chromedp.Action{
		chromedp.Navigate(category.URL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(2 * time.Second),
	}

	// Если FilterValue указан — выбираем производителя.
	// Если пустой — оставляем категорию без фильтра.
	if category.FilterValue != "" {
		manufacturerSelector := fmt.Sprintf(
			`input[type="checkbox"][value="%s"]`,
			category.FilterValue,
		)

		actions = append(actions,
			chromedp.WaitVisible(manufacturerSelector, chromedp.ByQuery),
			chromedp.Evaluate(
				fmt.Sprintf(`
					const el = document.querySelector('input[type="checkbox"][value="%s"]');
					if (el) {
						el.checked = true;
						el.dispatchEvent(new Event('change', { bubbles: true }));
					}
				`, category.FilterValue),
				nil,
			),
			chromedp.Sleep(2*time.Second),
		)
	}

	// Раскрываем все товары категории.
	actions = append(actions,
		chromedp.Evaluate(`
			const btn = document.querySelector(
				'span[title="Не запариваться и показать всё"]'
			);
			if (btn) btn.click();
		`, nil),

		chromedp.Sleep(2*time.Second),

		chromedp.Evaluate(js, &links),
	)

	err := chromedp.Run(p.ctx, actions...)
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
