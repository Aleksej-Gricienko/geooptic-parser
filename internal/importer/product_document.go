package importer

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"geooptic-parser-chromedp/models"
)

func ImportProductDocuments(
	db *sql.DB,
	productID int64,
	documents []models.ProductDocument,
) error {
	fmt.Printf("Documents: %d\n", len(documents))

	_, err := db.Exec(`
		DELETE FROM oc_product_document
		WHERE product_id = ?
	`, productID)

	if err != nil {
		return err
	}

	sortOrder := 0

	for _, document := range documents {
		fmt.Println("Import document:", document.Name)
		fmt.Println("Path:", document.Path)

		if document.Path == "" {
			fmt.Printf(
				"⚠️ Документ \"%s\" пропущен: путь пустой\n",
				document.Name,
			)
			continue
		}

		file := filepath.ToSlash(document.Path)

		_, err := db.Exec(`
			INSERT INTO oc_product_document
			(
				product_id,
				title,
				file,
				sort_order
			)
			VALUES
			(
				?,
				?,
				?,
				?
			)
		`,
			productID,
			document.Name,
			file,
			sortOrder,
		)

		if err != nil {
			return err
		}

		sortOrder++
	}

	return nil
}
