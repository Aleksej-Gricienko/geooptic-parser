package importer

import "database/sql"

const DefaultManufacturer = "Leica Geosystems"

func (db *Database) FindOrCreateManufacturer(name string) (int64, error) {
	var manufacturerID int64

	err := db.DB.QueryRow(`
		SELECT manufacturer_id
		FROM oc_manufacturer
		WHERE name = ?
		LIMIT 1
	`, name).Scan(&manufacturerID)

	if err == nil {
		return manufacturerID, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	result, err := db.DB.Exec(`
		INSERT INTO oc_manufacturer
		(
			name,
			image,
			sort_order
		)
		VALUES
		(
			?,
			'',
			0
		)
	`, name)

	if err != nil {
		return 0, err
	}

	manufacturerID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	_, err = db.DB.Exec(`
		INSERT INTO oc_manufacturer_to_store
		(
			manufacturer_id,
			store_id
		)
		VALUES
		(
			?,
			0
		)
	`, manufacturerID)

	if err != nil {
		return 0, err
	}

	return manufacturerID, nil
}
