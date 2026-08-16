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
		file         string
		folder       string
		manufacturer string
	}{
		// =========================
		// Leica Geosystems
		// =========================

		{"data/3d-scanners.json", "3d-scanners", "Leica Geosystems"},
		{"data/accessories.json", "accessories", "Leica Geosystems"},
		{"data/gnss-receivers.json", "gnss-receivers", "Leica Geosystems"},
		{"data/laser-levels.json", "laser-levels", "Leica Geosystems"},
		{"data/laser-rangefinders.json", "laser-rangefinders", "Leica Geosystems"},
		{"data/levels.json", "levels", "Leica Geosystems"},
		{"data/ndt-instruments.json", "ndt-instruments", "Leica Geosystems"},
		{"data/software.json", "software", "Leica Geosystems"},
		{"data/total-stations.json", "total-stations", "Leica Geosystems"},

		// =========================
		// GeoMax
		// =========================

		{"data/accessories-geomax.json", "accessories-geomax", "GeoMax"},
		{"data/gnss-receivers-geomax.json", "gnss-receivers-geomax", "GeoMax"},
		{"data/levels-geomax.json", "levels-geomax", "GeoMax"},
		{"data/ndt-instruments-geomax.json", "ndt-instruments-geomax", "GeoMax"},
		{"data/theodolites-geomax.json", "theodolites-geomax", "GeoMax"},
		{"data/total-stations-geomax.json", "total-stations-geomax", "GeoMax"},
	}

	for _, category := range categories {
		fmt.Printf(
			"\n========== IMPORT: %s ==========\n",
			category.folder,
		)

		err := db.ImportCategoryFile(
			category.file,
			category.folder,
			category.manufacturer,
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
