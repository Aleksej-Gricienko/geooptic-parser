package parser

import (
	"encoding/json"
	"os"

	"geooptic-parser-chromedp/models"
)

func SaveProducts(products []models.Product) error {
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	file, err := os.Create("data/products.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(products)
}
