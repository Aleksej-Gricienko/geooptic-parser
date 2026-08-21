package importer

import (
	"database/sql"
	"regexp"
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
	text = strings.ToLower(strings.TrimSpace(text))

	// Заменяем все символы, кроме латиницы, кириллицы и цифр, на "-"
	re := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	text = re.ReplaceAllString(text, "-")

	// Убираем "-" в начале и конце
	text = strings.Trim(text, "-")

	return text
}

func AddProductSEO(
	db *sql.DB,
	productID int64,
	keyword string,
) error {

	// Удаляем старые SEO-записи этого товара
	_, err := db.Exec(`
		DELETE FROM oc_seo_url
		WHERE store_id = 0
		  AND `+"`key`"+` = 'product_id'
		  AND `+"`value`"+` = ?
	`,
		productID,
	)

	if err != nil {
		return err
	}

	// Создаём SEO для английского и русского языков
	_, err = db.Exec(`
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
			(0, 1, 'product_id', ?, ?, 0),
			(0, 2, 'product_id', ?, ?, 0)
	`,
		productID,
		keyword,
		productID,
		keyword,
	)

	return err
}
