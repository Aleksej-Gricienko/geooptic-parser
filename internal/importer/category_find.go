package importer

func (db *Database) FindCategoryID(name string) (int64, error) {

	var id int64

	err := db.DB.QueryRow(`
SELECT category_id
FROM oc_category_description
WHERE name = ?
LIMIT 1
`, name).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
