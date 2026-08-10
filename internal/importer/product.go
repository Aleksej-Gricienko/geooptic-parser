package importer

import (
	"database/sql"
	"time"

	"geooptic-parser-chromedp/models"
)

type ProductResult struct {
	ID int64
}

func ImportProduct(
	db *sql.DB,
	product models.Product,
	manufacturerID int64,
	images []string,
) (*ProductResult, error) {

	now := time.Now()
	image := ""

	if len(images) > 0 {
		image = images[0]
	}

	result, err := db.Exec(`
INSERT INTO oc_product
(
	model,
	quantity,
	stock_status_id,
	image,
	manufacturer_id,
	shipping,
	price,
	date_available,
	status,
	date_added,
	date_modified
)
VALUES
(
	?,
	1,
	7,
	?,
	?,
	1,
	0,
	CURDATE(),
	1,
	?,
	?
)
`,
		product.Name,
		image,
		manufacturerID,
		now,
		now,
	)

	if err != nil {
		return nil, err
	}

	productID, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}
	metaDescription := product.Name
	if len(metaDescription) > 250 {
		metaDescription = metaDescription[:250] + "..."
	}

	_, err = db.Exec(`
INSERT INTO oc_product_description
(
	product_id,
	language_id,
	name,
	description,
	tag,
	meta_title,
	meta_description,
	meta_keyword
)
VALUES
(
	?,
	1,
	?,
	?,
	'',
	?,
	?,
	''
)
`,
		productID,
		product.Name,
		product.DescriptionHTML,
		product.Name,    // meta_title
		metaDescription, // meta_description
	)

	if err != nil {
		return nil, err
	}

	return &ProductResult{
		ID: productID,
	}, nil
}
