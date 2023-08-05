package manejador

import (
	"net/http"
)

type Manejador struct {
	Archivos      []Archivo
	Cliente       http.Client
	Configuracion Configuracion
	Informe       [][]string
	Errores       bool
}

type Archivo struct {
	Extension string // Minúsculas
	Nombre    string
	Ruta      string
}

type Configuracion struct {
	Dialogos           bool
	Comando            string
	ExtensionesValidas []string // Minúsculas
	Token              string
	RutaOrigen         string
	RutaDestino        string
	RutaConfiguracion  string
	RutaImageMagick    string
}

type RespuestaApi struct {
	Results []struct {
		Box struct {
			Xmin int `json:"xmin"`
			Ymin int `json:"ymin"`
			Xmax int `json:"xmax"`
			Ymax int `json:"ymax"`
		} `json:"box"`
	} `json:"results"`
}
