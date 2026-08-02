package internal

import (
	"fmt"
	"geooptic-parser-chromedp/models"
	"os"
	"path/filepath"
	"strings"
)

const DocumentRoot = "/var/www/leica-cms/upload/documents/catalog/Categories/total-stations"

func DownloadDocuments(products []models.Product) error {

	if err := os.MkdirAll(DocumentRoot, 0755); err != nil {
		return err
	}

	for _, product := range products {

		productDir := filepath.Join(DocumentRoot, createSlug(product.Name))

		if err := os.MkdirAll(productDir, 0755); err != nil {
			fmt.Println("Ошибка создания папки:", err)
			continue
		}

		if len(product.Documents) == 0 {
			fmt.Println("Нет PDF:", product.Name)
			continue
		}

		for i, doc := range product.Documents {

			filename := doc.Name + ".pdf"

			if !strings.HasSuffix(strings.ToLower(filename), ".pdf") {
				filename += ".pdf"
			}

			fullPath := filepath.Join(productDir, filename)

			fmt.Println("Скачиваю PDF:", filename)

			err := downloadFile(doc.URL, fullPath)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			product.Documents[i].Path = fullPath
			fmt.Println("Сохранено:", fullPath)
		}
	}

	return nil
}
