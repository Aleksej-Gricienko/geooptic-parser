package importer

import (
	"fmt"
)

func (db *Database) ImportCategoryFile(
	path string,
	categoryFolder string,
) error {
	info, ok := categoryMap[categoryFolder]
	if !ok {
		return fmt.Errorf("неизвестная категория: %s", categoryFolder)
	}

	categoryID, err := db.FindCategoryID(info.Name)
	if err != nil {
		return fmt.Errorf(
			"категория %q не найдена в БД: %w",
			info.Name,
			err,
		)
	}

	fmt.Println("Category name:", info.Name)
	fmt.Println("Category folder:", info.Folder)
	fmt.Println("Category ID:", categoryID)

	products, err := LoadProducts(path)
	if err != nil {
		return err
	}

	fmt.Printf("Products: %d\n", len(products))

	if len(products) == 0 {
		return fmt.Errorf("в файле %s нет товаров", path)
	}

	manufacturerID, err := db.FindOrCreateManufacturer(DefaultManufacturer)
	if err != nil {
		return err
	}

	for i, product := range products {
		fmt.Printf(
			"\n[%d/%d] %s\n",
			i+1,
			len(products),
			product.Name,
		)

		fmt.Printf(
			"Documents: %d\n",
			len(product.Documents),
		)

		// Ищем уже существующие картинки.
		// Документы здесь вообще не участвуют.
		folder := ProductImageFolder(product.Name)

		fmt.Println("Product folder:", folder)

		images, err := FindImages(info.Folder, folder)
		if err != nil {
			return fmt.Errorf(
				"ошибка поиска изображений товара %q: %w",
				product.Name,
				err,
			)
		}

		fmt.Printf("Images: %d\n", len(images))

		result, err := ImportProduct(
			db.DB,
			product,
			manufacturerID,
			images,
		)
		if err != nil {
			return fmt.Errorf(
				"ошибка импорта товара %q: %w",
				product.Name,
				err,
			)
		}
		fmt.Printf(
			"Characteristics: %d\n",
			len(product.Characteristics),
		)

		for name, value := range product.Characteristics {
			attributeID, err := db.FindOrCreateAttribute(name)
			if err != nil {
				return fmt.Errorf(
					"ошибка поиска/создания атрибута %q: %w",
					name,
					err,
				)
			}

			err = AddProductAttribute(
				db.DB,
				result.ID,
				attributeID,
				value,
			)
			if err != nil {
				return fmt.Errorf(
					"ошибка импорта атрибута %q товара %q: %w",
					name,
					product.Name,
					err,
				)
			}
		}

		err = ImportProductImages(
			db.DB,
			result.ID,
			images,
		)
		if err != nil {
			return fmt.Errorf(
				"ошибка импорта изображений товара %q: %w",
				product.Name,
				err,
			)
		}

		// А документы — совершенно отдельно.
		err = ImportProductDocuments(
			db.DB,
			result.ID,
			product.Documents,
		)
		if err != nil {
			return fmt.Errorf(
				"ошибка импорта документов товара %q: %w",
				product.Name,
				err,
			)
		}

		err = AddProductSEO(
			db.DB,
			result.ID,
			ProductSEOKeyword(product.Name),
		)
		if err != nil {
			return err
		}

		err = db.LinkProductToCategory(
			result.ID,
			categoryID,
		)
		if err != nil {
			return err
		}

		err = db.LinkProductToStore(result.ID)
		if err != nil {
			return err
		}

		fmt.Printf(
			"✅ Imported: %s (ID=%d)\n",
			product.Name,
			result.ID,
		)
	}

	return nil
}
