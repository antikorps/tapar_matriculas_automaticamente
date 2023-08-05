package manejador

import (
	"path/filepath"
	"strings"
)

func nombreDeArchivo(ruta string) string {
	extension := filepath.Ext(ruta)
	return strings.TrimSuffix(filepath.Base(ruta), extension)
}

func contiene(coleccion []string, elemento string) bool {
	for _, v := range coleccion {
		if v == elemento {
			return true
		}
	}
	return false
}
