package main

import (
	"fmt"
	"log"

	"geooptic-parser-chromedp/internal/importer"
)

func main() {
	db, err := importer.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.DB.Close()

	categories := []struct {
		file   string
		folder string
	}{
		{"data/3d-scanners.json", "3d-scanners"},
		{"data/accessories.json", "accessories"},
		{"data/gnss-receivers.json", "gnss-receivers"},
		{"data/laser-levels.json", "laser-levels"},
		{"data/laser-rangefinders.json", "laser-rangefinders"},
		{"data/levels.json", "levels"},
		{"data/ndt-instruments.json", "ndt-instruments"},
		{"data/software.json", "software"},
		{"data/total-stations.json", "total-stations"},
	}

	for _, category := range categories {
		fmt.Printf(
			"\n========== IMPORT: %s ==========\n",
			category.folder,
		)

		err := db.ImportCategoryFile(
			category.file,
			category.folder,
		)
		if err != nil {
			log.Fatalf(
				"ошибка импорта категории %s: %v",
				category.folder,
				err,
			)
		}
	}

	fmt.Println("\n🎉 Импорт всех категорий завершён.")
}
