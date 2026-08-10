package importer

import "database/sql"

func AddProductAttribute(
	db *sql.DB,
	productID int64,
	attributeID int64,
	value string,
) error {

	_, err := db.Exec(`
		INSERT INTO oc_product_attribute
		(
			product_id,
			attribute_id,
			language_id,
			text
		)
		VALUES
		(
			?,
			?,
			1,
			?
		)
	`,
		productID,
		attributeID,
		value,
	)

	return err
}
