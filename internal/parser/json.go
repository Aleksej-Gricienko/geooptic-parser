package parser

import (
	"encoding/json"
	"os"

	"geooptic-parser-chromedp/models"
)

func SaveProducts(products []models.Product, categorySlug string) error {
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	file, err := os.Create("data/" + categorySlug + ".json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(products)
}
