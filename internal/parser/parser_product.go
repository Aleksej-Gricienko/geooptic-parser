package parser

import (
	"fmt"
	"geooptic-parser-chromedp/models"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func (p *Parser) ParseProduct(url string) (models.Product, error) {

	// -----------------------
	// Первая загрузка страницы
	// -----------------------

	html, err := p.getHTML(url)
	if err != nil {
		return models.Product{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return models.Product{}, fmt.Errorf("ошибка разбора HTML: %w", err)
	}

	product := models.Product{
		URL:             url,
		Characteristics: make(map[string]string),
	}

	product.Name = parseName(doc)
	product.Manufacturer = parseManufacturer(doc)
	product.DescriptionHTML, product.DescriptionText = parseDescription(doc)
	product.Characteristics = parseCharacteristics(doc)
	product.Images = parseImages(doc)

	// -----------------------
	// Вторая загрузка страницы
	// Только для PDF
	// -----------------------
	fmt.Println("=== ПЕРЕХОД К ФАЙЛАМ ===")
	fmt.Println("Товар:", product.Name)
	filesHTML, err := p.getHTMLWithFiles(url)
	if err != nil {
		return models.Product{}, err
	}

	filesDoc, err := goquery.NewDocumentFromReader(strings.NewReader(filesHTML))
	if err != nil {
		return models.Product{}, fmt.Errorf("ошибка разбора HTML с PDF: %w", err)
	}

	filesDoc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		text := strings.TrimSpace(s.Text())

		if strings.Contains(strings.ToLower(href), ".pdf") {
			fmt.Println("TEXT =", text)
			fmt.Println("HREF =", href)
		}
	})

	product.Documents = parseDocuments(filesDoc)
	fmt.Println("Документов найдено:", len(product.Documents))
	fmt.Printf("%+v\n", product.Documents)

	// -----------------------
	// Отладка
	// -----------------------

	fmt.Println("URL:", product.URL)

	fmt.Println("\nОписание:")
	fmt.Println(product.DescriptionText)

	fmt.Println("\nХарактеристики:")
	for key, value := range product.Characteristics {
		fmt.Printf("%s: %s\n", key, value)
	}

	fmt.Println("\nPDF:")
	for _, pdf := range product.Documents {
		fmt.Printf("%s -> %s\n", pdf.Name, pdf.URL)
	}

	return product, nil
}

func (p *Parser) getHTML(url string) (string, error) {

	var html string

	err := chromedp.Run(
		p.ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)

	if err != nil {
		return "", fmt.Errorf("ошибка загрузки страницы %s: %w", url, err)
	}

	return html, nil
}

func (p *Parser) getHTMLWithFiles(url string) (string, error) {
	var html string
	var hasFilesTab bool

	fmt.Println("Открываю страницу файлов:", url)

	// Сначала просто открываем страницу.
	err := chromedp.Run(
		p.ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),

		chromedp.Evaluate(`
			(() => {
				return [...document.querySelectorAll("span")]
					.some(el => el.textContent.trim() === "Файлы");
			})()
		`, &hasFilesTab),
	)

	if err != nil {
		return "", fmt.Errorf(
			"ошибка загрузки страницы %s: %w",
			url,
			err,
		)
	}

	// Вкладки «Файлы» нет — у товара просто нет документов.
	if !hasFilesTab {
		fmt.Println("Вкладка «Файлы» отсутствует — документов нет")
		return "", nil
	}

	fmt.Println("Вкладка «Файлы» найдена")

	// Открываем вкладку.
	err = chromedp.Run(
		p.ctx,

		chromedp.Click(
			`//span[text()="Файлы"]`,
			chromedp.BySearch,
		),

		chromedp.Sleep(1*time.Second),

		chromedp.OuterHTML("html", &html, chromedp.ByQuery),
	)

	if err != nil {
		return "", fmt.Errorf(
			"ошибка открытия вкладки файлов %s: %w",
			url,
			err,
		)
	}

	fmt.Println("HTML файлов получен, размер:", len(html))

	return html, nil
}
func parseName(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("h1").First().Text())
}

func parseDescription(doc *goquery.Document) (string, string) {

	desc := doc.Find(".description").First()

	if strings.TrimSpace(desc.Text()) != "" {

		html, _ := desc.Html()

		return strings.TrimSpace(html), strings.TrimSpace(desc.Text())
	}

	short := doc.Find(".shortDesc").First()

	html, _ := short.Html()

	return strings.TrimSpace(html), strings.TrimSpace(short.Text())
}

func parseCharacteristics(doc *goquery.Document) map[string]string {

	result := make(map[string]string)

	doc.Find("table.params tbody tr").Each(func(i int, s *goquery.Selection) {

		td := s.Find("td")

		if td.Length() == 3 {

			key := strings.TrimSpace(td.Eq(1).Text())
			value := strings.TrimSpace(td.Eq(2).Text())

			if key != "" && value != "" {
				result[key] = value
			}

		}

		if td.Length() == 2 {

			key := strings.TrimSpace(td.Eq(0).Text())
			value := strings.TrimSpace(td.Eq(1).Text())

			if key != "" && value != "" {
				result[key] = value
			}

		}

	})

	return result
}
func parseManufacturer(doc *goquery.Document) string {
	var manufacturer string

	doc.Find("*").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := strings.TrimSpace(s.Text())

		if text != "Производитель" {
			return true
		}

		parent := s.Parent()

		if parent.Length() == 0 {
			return true
		}

		parent.Find("a").EachWithBreak(func(i int, a *goquery.Selection) bool {
			value := strings.TrimSpace(a.Text())

			if value != "" && value != "Производитель" {
				manufacturer = value
				return false
			}

			return true
		})

		return manufacturer == ""
	})

	return manufacturer
}

func parseImages(doc *goquery.Document) []string {
	var images []string

	doc.Find(".imgs .img img").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok {
			return
		}

		if strings.HasSuffix(strings.ToLower(src), "-1.jpg") {
			src = src[:len(src)-len("-1.jpg")] + "-0.jpg"
		}

		images = append(images, src)
	})

	return images
}

func parseDocuments(doc *goquery.Document) []models.ProductDocument {

	var docs []models.ProductDocument
	seen := make(map[string]bool)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {

		href, ok := s.Attr("href")
		if !ok {
			return
		}

		lower := strings.ToLower(href)
		if strings.EqualFold(
			filepath.Base(href),
			"geooptic_sertificate_2022.pdf",
		) {
			return
		}

		if !strings.HasSuffix(lower, ".pdf") &&
			!strings.HasSuffix(lower, ".ppt") &&
			!strings.HasSuffix(lower, ".pptx") {
			return
		}

		if seen[href] {
			return
		}

		seen[href] = true

		name := strings.TrimSpace(s.Text())
		if name == "" || name == "Скачать" {
			base := filepath.Base(href)
			ext := filepath.Ext(base)
			name = strings.TrimSuffix(base, ext)
		}
		docs = append(docs, models.ProductDocument{
			Name: name,
			URL:  href,
		})
	})
	fmt.Println("=== Documents ===")

	for _, d := range docs {
		fmt.Println(d.Name)
		fmt.Println(d.URL)
		fmt.Println()
	}
	return docs
}
