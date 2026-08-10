package importer

import "database/sql"

func ImportProductImages(
	db *sql.DB,
	productID int64,
	images []string,
) error {

	if len(images) <= 1 {
		return nil
	}

	sortOrder := 0

	for _, image := range images[1:] {

		_, err := db.Exec(`
			INSERT INTO oc_product_image
			(
				product_id,
				image,
				sort_order
			)
			VALUES
			(
				?,
				?,
				?
			)
		`,
			productID,
			image,
			sortOrder,
		)

		if err != nil {
			return err
		}

		sortOrder++
	}

	return nil
}
