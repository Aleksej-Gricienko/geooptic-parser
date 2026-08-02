package models

type ProductDocument struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Path string `json:"path"`
}

type Product struct {
	Name            string            `json:"name"`
	URL             string            `json:"url"`
	DescriptionHTML string            `json:"description_html"`
	DescriptionText string            `json:"description_text"`
	Characteristics map[string]string `json:"characteristics"`
	Images          []string          `json:"images"`
	Documents       []ProductDocument `json:"documents"`
}
