package importer

import (
	"encoding/json"
	"os"

	"geooptic-parser-chromedp/models"
)

func LoadProducts(path string) ([]models.Product, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var products []models.Product

	err = json.Unmarshal(data, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}
