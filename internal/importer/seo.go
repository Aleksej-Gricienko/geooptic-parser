package importer

import (
	"database/sql"
	"strings"
)

func ProductSEOKeyword(name string) string {

	name = strings.TrimSpace(name)

	index := strings.Index(strings.ToLower(name), "leica")

	if index == -1 {
		return Slugify(name)
	}

	return Slugify(name[index:])
}

func Slugify(text string) string {

	text = strings.ToLower(text)

	text = strings.ReplaceAll(text, " ", "-")

	return text
}

func AddProductSEO(
	db *sql.DB,
	productID int64,
	keyword string,
) error {

	_, err := db.Exec(`
		INSERT INTO oc_seo_url
		(
			store_id,
			language_id,
			`+"`key`"+`,
			value,
			keyword,
			sort_order
		)
		VALUES
		(
			0,
			1,
			'product_id',
			?,
			?,
			0
		)
	`,
		productID,
		keyword,
	)

	return err
}
