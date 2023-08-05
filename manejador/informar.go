package manejador

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ncruces/zenity"
)

func (m *Manejador) Informar() {
	nombreInforme := fmt.Sprintf("informe_procesamiento_%d.csv", time.Now().Unix())
	rutaInforme := filepath.Join(m.Configuracion.RutaDestino, nombreInforme)

	archivoCSV, archivoCSVError := os.Create(rutaInforme)
	if archivoCSVError != nil {
		mensajeError := fmt.Sprintf("el procesamiento de las imágenes ha terminado, pero no podido crearse el informe de la ejecución: %v", archivoCSVError)
		if m.Configuracion.Dialogos {
			zenity.Info(mensajeError,
				zenity.Title("ERROR CRÍTICO"),
				zenity.Width(600),
				zenity.ErrorIcon)
		}
		log.Fatalln(mensajeError)
	}
	defer archivoCSV.Close()

	manejadorCSV := csv.NewWriter(archivoCSV)
	errorEscritura := manejadorCSV.WriteAll(m.Informe)
	if errorEscritura != nil {
		mensajeError := fmt.Sprintf("el procesamiento de las imágenes ha terminado, pero no podido escribirse el informe de la ejecución: %v", errorEscritura)
		if m.Configuracion.Dialogos {
			zenity.Info(mensajeError,
				zenity.Title("ERROR CRÍTICO"),
				zenity.Width(600),
				zenity.ErrorIcon)
		}
		log.Fatalln(mensajeError)
	}

	if m.Errores {
		mensajeError := fmt.Sprintf("ATENCIÓN: se han producido errores durante el procesamiento de las imágenes.\nConsulta el informe generado en:\n%v", rutaInforme)
		if m.Configuracion.Dialogos {
			zenity.Info(mensajeError,
				zenity.Title("FIN"),
				zenity.Width(600),
				zenity.ErrorIcon)
		}
		log.Fatalln(mensajeError)
	}

	if m.Configuracion.Dialogos {
		mensajeFin := fmt.Sprintf("Procesamiento realizado con éxito.\nConsulta el informe generado en:\n %v", rutaInforme)
		if m.Configuracion.Dialogos {
			zenity.Info(mensajeFin,
				zenity.Title("FIN"),
				zenity.Width(600),
				zenity.InfoIcon)
		}
	}
}
