package main

import (
	"flag"
	"tapar_matriculas_automaticamente/manejador"
)

func main() {
	var origen, destino string
	flag.StringVar(&origen, "origen", "", "ruta completa a la imagen o al directorio a analizar")
	flag.StringVar(&destino, "destino", "", "ruta completa del directorio en el que se guardarán las nuevas imágenes")
	flag.Parse()

	m := manejador.Crear(origen, destino)
	m.Procesar()
	m.Informar()

}
