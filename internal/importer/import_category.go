package importer

import (
	"fmt"
	"path/filepath"
)

func (db *Database) ImportCategoryFile(path string) error {

	info := CategoryInfoFromFile(path)

	categoryID, err := db.FindCategoryID(info.Name)
	fmt.Println("Category name:", info.Name)
	fmt.Println("Category folder:", info.Folder)
	if err != nil {
		return err
	}

	products, err := LoadProducts(path)
	if err != nil {
		return err
	}

	fmt.Printf("Category: %s (ID=%d)\n", info.Name, categoryID)
	fmt.Printf("Products: %d\n", len(products))

	if len(products) == 0 {
		return fmt.Errorf("в файле %s нет товаров", path)
	}

	// Пока импортируем только первый товар

	for i, product := range products {
		fmt.Printf("[%d/%d] %s\n", i+1, len(products), product.Name)
		fmt.Printf("\n=========================\n")
		fmt.Printf("Import: %s\n", product.Name)
		fmt.Println()
		fmt.Printf("JSON товаров: %d\n", len(products))

		if len(product.Documents) == 0 {
			fmt.Printf("⚠️ Пропущен %s: нет документов\n", product.Name)
			continue
		}

		folder := filepath.Base(filepath.Dir(product.Documents[0].Path))

		fmt.Println("Product folder:", folder)

		images, err := FindImages(info.Folder, folder)
		if err != nil {
			return err
		}

		manufacturerID, err := db.FindOrCreateManufacturer(DefaultManufacturer)
		if err != nil {
			return err
		}

		result, err := ImportProduct(
			db.DB,
			product,
			manufacturerID,
			images,
		)
		if err != nil {
			return err
		}

		err = ImportProductImages(db.DB, result.ID, images)
		if err != nil {
			return err
		}

		for name, value := range product.Characteristics {

			attributeID, err := db.FindOrCreateAttribute(name)
			if err != nil {
				return err
			}

			err = AddProductAttribute(
				db.DB,
				result.ID,
				attributeID,
				value,
			)
			if err != nil {
				return err
			}
		}

		err = ImportProductDocuments(
			db.DB,
			result.ID,
			product.Documents,
		)
		if err != nil {
			return err
		}

		keyword := ProductSEOKeyword(product.Name)

		err = AddProductSEO(
			db.DB,
			result.ID,
			keyword,
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

		fmt.Printf("✅ Imported: %s (ID=%d)\n", product.Name, result.ID)
	}

	return nil
}
