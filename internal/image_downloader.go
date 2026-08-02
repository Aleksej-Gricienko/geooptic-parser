package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"geooptic-parser-chromedp/models"
)

const ImageRoot = "/var/www/leica-cms/upload/image/catalog/Categories/ndt-instruments"

func createSlug(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	return name
}
func downloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
func DownloadImages(products []models.Product) error {

	if err := os.MkdirAll(ImageRoot, 0755); err != nil {
		return err
	}

	for _, product := range products {
		productDir := filepath.Join(ImageRoot, createSlug(product.Name))

		if err := os.MkdirAll(productDir, 0755); err != nil {
			fmt.Println("Ошибка создания папки:", err)
			continue
		}

		if len(product.Images) == 0 {
			fmt.Println("Нет изображения:", product.Name)
			continue
		}

		for i, imageURL := range product.Images {

			filename := fmt.Sprintf("%d.jpg", i+1)

			fullPath := filepath.Join(productDir, filename)

			fmt.Println("Скачиваю:", filename)

			err := downloadFile(imageURL, fullPath)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
		}
	}
	return nil
}
