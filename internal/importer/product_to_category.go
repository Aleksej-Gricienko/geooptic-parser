package importer

func (db *Database) LinkProductToCategory(
	productID int64,
	categoryID int64,
) error {

	_, err := db.DB.Exec(`
INSERT INTO oc_product_to_category
(
	product_id,
	category_id
)
VALUES
(
	?,
	?
)
`,
		productID,
		categoryID,
	)

	return err
}
