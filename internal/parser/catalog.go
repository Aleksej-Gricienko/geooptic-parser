package parser

import (
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

const catalogURL = "https://www.geooptic.ru/catalog/lazernye-urovni"

func (p *Parser) GetProductLinks() ([]string, error) {

	var links []string

	js := `
	Array.from(document.querySelectorAll(".productList .productItem > a.img"))
    .map(a => a.href);
	`

	err := chromedp.Run(
		p.ctx,

		chromedp.Navigate(catalogURL),

		chromedp.WaitVisible(`input[type="checkbox"][value="3"]`),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`
const el = document.querySelector('input[type="checkbox"][value="3"]');
el.checked = true;
el.dispatchEvent(new Event('change', { bubbles: true }));
`, nil),

		chromedp.Sleep(2*time.Second),
		// Раскрываем все товары категории
		chromedp.Evaluate(`
		const btn = document.querySelector('span[title="Не запариваться и показать всё"]');
		if (btn) btn.click();
		`, nil),

		chromedp.Sleep(2*time.Second),

		chromedp.Evaluate(js, &links))

	if err != nil {
		return nil, fmt.Errorf("не удалось получить ссылки: %w", err)
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
