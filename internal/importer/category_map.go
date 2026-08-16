package importer

import "path/filepath"

type CategoryInfo struct {
	Name   string
	Folder string
}

var categoryMap = map[string]CategoryInfo{
	"3d-scanners": {
		Name:   "3D Сканеры",
		Folder: "total-stations/3d-scanners",
	},
	"accessories": {
		Name:   "Аксессуары",
		Folder: "total-stations/accessories",
	},
	"gnss-receivers": {
		Name:   "GNSS-приемники",
		Folder: "total-stations/gnss-receivers",
	},
	"levels": {
		Name:   "Нивелиры",
		Folder: "total-stations/levels",
	},
	"laser-levels": {
		Name:   "Лазерные уровни",
		Folder: "total-stations/laser-levels",
	},
	"laser-rangefinders": {
		Name:   "Лазерные дальномеры",
		Folder: "total-stations/laser-rangefinders",
	},
	"ndt-instruments": {
		Name:   "Приборы НК",
		Folder: "total-stations/ndt-instruments",
	},
	"software": {
		Name:   "Программное обеспечение",
		Folder: "total-stations/software",
	},
	"total-stations": {
		Name:   "Тахеометры",
		Folder: "total-stations/total-stations",
	},

	// GeoMax
	"accessories-geomax": {
		Name:   "Аксессуары",
		Folder: "total-stations/accessories-geomax",
	},
	"gnss-receivers-geomax": {
		Name:   "GNSS-приемники",
		Folder: "total-stations/gnss-receivers-geomax",
	},
	"levels-geomax": {
		Name:   "Нивелиры",
		Folder: "total-stations/levels-geomax",
	},
	"ndt-instruments-geomax": {
		Name:   "Приборы НК",
		Folder: "total-stations/ndt-instruments-geomax",
	},
	"theodolites-geomax": {
		Name:   "Теодолиты",
		Folder: "total-stations/theodolites-geomax",
	},
	"total-stations-geomax": {
		Name:   "Тахеометры",
		Folder: "total-stations/total-stations-geomax",
	},
}

func CategoryInfoFromFile(path string) CategoryInfo {
	name := filepath.Base(path)
	ext := filepath.Ext(name)

	name = name[:len(name)-len(ext)]

	return categoryMap[name]
}
