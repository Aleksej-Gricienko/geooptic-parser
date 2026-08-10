package importer

func (db *Database) LinkProductToStore(
	productID int64,
) error {

	_, err := db.DB.Exec(`
		INSERT INTO oc_product_to_store
		(
			product_id,
			store_id
		)
		VALUES
		(
			?,
			0
		)
	`, productID)

	return err
}
