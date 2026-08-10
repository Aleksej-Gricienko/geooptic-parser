package importer

import "database/sql"

const AttributeGroupID = 7

func (db *Database) FindOrCreateAttribute(name string) (int64, error) {

	var attributeID int64

	err := db.DB.QueryRow(`
		SELECT attribute_id
		FROM oc_attribute_description
		WHERE name = ?
		LIMIT 1
	`, name).Scan(&attributeID)

	if err == nil {
		return attributeID, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	result, err := db.DB.Exec(`
		INSERT INTO oc_attribute
		(
			attribute_group_id,
			sort_order
		)
		VALUES
		(
			?,
			0
		)
	`, AttributeGroupID)

	if err != nil {
		return 0, err
	}

	attributeID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	_, err = db.DB.Exec(`
		INSERT INTO oc_attribute_description
		(
			attribute_id,
			language_id,
			name
		)
		VALUES
		(
			?,
			1,
			?
		)
	`, attributeID, name)

	if err != nil {
		return 0, err
	}

	return attributeID, nil
}
