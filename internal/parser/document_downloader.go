package parser

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"geooptic-parser-chromedp/models"
)

const DocumentRoot = "/var/www/leica-cms/upload/documents/catalog/Categories"

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"\"", "",
		"*", "",
		"?", "",
		"<", "",
		">", "",
		"|", "",
	)

	return replacer.Replace(strings.TrimSpace(name))
}

func fileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()

	buffer := make([]byte, 1024*1024)

	for {
		n, err := file.Read(buffer)

		if n > 0 {
			if _, writeErr := hash.Write(buffer[:n]); writeErr != nil {
				return "", writeErr
			}
		}

		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			return "", err
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func DownloadDocuments(products []models.Product, category string) error {
	categorySlug := createSlug(category)

	categoryDir := filepath.Join(DocumentRoot, categorySlug)

	if err := os.MkdirAll(categoryDir, 0755); err != nil {
		return err
	}

	// hash -> относительный путь к уже существующему файлу
	hashToPath := make(map[string]string)

	// URL -> относительный путь.
	// Это позволяет не скачивать один и тот же URL несколько раз
	// в рамках одного запуска.
	urlToPath := make(map[string]string)

	// Индексируем уже существующие файлы категории.
	entries, err := os.ReadDir(categoryDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(categoryDir, entry.Name())

		hash, err := fileHash(fullPath)
		if err != nil {
			fmt.Println("Не удалось вычислить hash:", fullPath, err)
			continue
		}

		relativePath := filepath.ToSlash(
			filepath.Join(
				"documents",
				"catalog",
				"Categories",
				categorySlug,
				entry.Name(),
			),
		)

		hashToPath[hash] = relativePath
	}

	fmt.Printf(
		"Категория: %s | существующих документов: %d\n",
		categorySlug,
		len(hashToPath),
	)

	for productIndex := range products {
		product := &products[productIndex]

		if len(product.Documents) == 0 {
			fmt.Println("Нет документов:", product.Name)
			continue
		}

		fmt.Printf(
			"\nТовар: %s | документов: %d\n",
			product.Name,
			len(product.Documents),
		)

		for documentIndex := range product.Documents {
			document := &product.Documents[documentIndex]

			// -----------------------------------
			// 1. Такой URL уже скачивали
			// -----------------------------------

			if path, exists := urlToPath[document.URL]; exists {
				document.Path = path

				fmt.Println("Уже скачан:", document.Name)
				fmt.Println("Path:", path)

				continue
			}

			// -----------------------------------
			// 2. Определяем имя файла
			// -----------------------------------

			ext := strings.ToLower(filepath.Ext(document.URL))

			if ext == "" {
				ext = ".pdf"
			}

			filename := sanitizeFilename(document.Name)

			if filename == "" {
				filename = "document"
			}

			if !strings.HasSuffix(strings.ToLower(filename), ext) {
				filename += ext
			}

			// -----------------------------------
			// 3. Временный файл
			// -----------------------------------

			tempFile, err := os.CreateTemp(
				categoryDir,
				".document-*",
			)

			if err != nil {
				fmt.Println("Ошибка создания временного файла:", err)
				continue
			}

			tempPath := tempFile.Name()

			if err := tempFile.Close(); err != nil {
				os.Remove(tempPath)
				fmt.Println("Ошибка закрытия временного файла:", err)
				continue
			}

			fmt.Println("Скачиваю:", document.Name)

			err = downloadFile(document.URL, tempPath)

			if err != nil {
				os.Remove(tempPath)

				fmt.Println("Ошибка скачивания:", err)
				continue
			}
			if err := os.Chmod(tempPath, 0644); err != nil {
				fmt.Println("Ошибка установки прав:", err)
				os.Remove(tempPath)
				continue
			}
			// -----------------------------------
			// 4. Вычисляем SHA-256
			// -----------------------------------

			hash, err := fileHash(tempPath)

			if err != nil {
				os.Remove(tempPath)

				fmt.Println("Ошибка вычисления hash:", err)
				continue
			}

			// -----------------------------------
			// 5. Такой файл уже существует
			// -----------------------------------

			if existingPath, exists := hashToPath[hash]; exists {
				os.Remove(tempPath)

				urlToPath[document.URL] = existingPath
				document.Path = existingPath

				fmt.Println("Дубликат:", document.Name)
				fmt.Println("Использую:", existingPath)

				continue
			}

			// -----------------------------------
			// 6. Файл новый
			// -----------------------------------

			finalPath := filepath.Join(categoryDir, filename)

			// Если файл с таким именем уже существует,
			// но содержимое другое — не перезаписываем его.
			if _, err := os.Stat(finalPath); err == nil {
				ext := filepath.Ext(filename)
				nameWithoutExt := strings.TrimSuffix(filename, ext)

				filename = fmt.Sprintf(
					"%s-%s%s",
					nameWithoutExt,
					hash[:12],
					ext,
				)

				finalPath = filepath.Join(categoryDir, filename)
			}

			if err := os.Rename(tempPath, finalPath); err != nil {
				os.Remove(tempPath)

				fmt.Println("Ошибка сохранения:", err)
				continue
			}

			relativePath := filepath.ToSlash(
				filepath.Join(
					"documents",
					"catalog",
					"Categories",
					categorySlug,
					filename,
				),
			)

			hashToPath[hash] = relativePath
			urlToPath[document.URL] = relativePath

			document.Path = relativePath

			fmt.Println("Сохранено:", relativePath)
			fmt.Println("SHA-256:", hash)
		}
	}

	return nil
}
