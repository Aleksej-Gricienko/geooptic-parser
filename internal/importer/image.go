package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const OpenCartImageDir = "/var/www/leica-cms/upload/image/catalog/Categories"

func FindImages(categoryFolder, productName string) ([]string, error) {

	dir := filepath.Join(
		OpenCartImageDir,
		categoryFolder,
		productName,
	)
	fmt.Println("=== FindImages ===")
	fmt.Println("Directory:", dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var images []string

	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))

		switch ext {
		case ".jpg", ".jpeg", ".png", ".webp":
			images = append(
				images,
				filepath.ToSlash(filepath.Join(
					"catalog",
					"Categories",
					categoryFolder,
					productName,
					entry.Name(),
				)),
			)
		}
	}

	sort.Strings(images)
	fmt.Println("Directory:", dir)
	fmt.Printf("Entries: %d\n", len(entries))

	for _, entry := range entries {
		fmt.Println(entry.Name())
	}

	fmt.Printf("Found images: %d\n", len(images))
	return images, nil
}
func ProductImageFolder(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	name = strings.ReplaceAll(name, "\"", "-")

	return name
}
