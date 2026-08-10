package importer

import "path/filepath"

type CategoryInfo struct {
	Name   string
	Folder string
}

var categoryMap = map[string]CategoryInfo{
	"3d-scanners": {
		Name:   "3D Сканеры",
		Folder: "3D-scanners",
	},
	"accessories": {
		Name:   "Аксессуары",
		Folder: "accessories",
	},
	"gnss-receivers": {
		Name:   "GNSS-приемники",
		Folder: "gnss-receivers",
	},
	"levels": {
		Name:   "Нивелиры",
		Folder: "levels",
	},
	"laser-levels": {
		Name:   "Лазерные уровни",
		Folder: "laser-levels",
	},
	"laser-rangefinders": {
		Name:   "Лазерные дальномеры",
		Folder: "laser-rangefinders",
	},
	"ndt-instruments": {
		Name:   "Приборы НК",
		Folder: "ndt-instruments",
	},
	"software": {
		Name:   "Программное обеспечение",
		Folder: "software",
	},
	"total-stations": {
		Name:   "Тахеометры",
		Folder: "total-stations",
	},
}

func CategoryInfoFromFile(path string) CategoryInfo {
	name := filepath.Base(path)
	ext := filepath.Ext(name)

	name = name[:len(name)-len(ext)]

	return categoryMap[name]
}
