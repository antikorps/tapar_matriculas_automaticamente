package manejador

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ncruces/zenity"
)

func peticionAPI(cliente http.Client, rutaOrigen, rutaDestino, token string) (RespuestaApi, error) {
	var respuestaApi RespuestaApi
	form := new(bytes.Buffer)
	manejadorForm := multipart.NewWriter(form)
	manejadorFile, manejadorFileError := manejadorForm.CreateFormFile("upload", filepath.Base(rutaOrigen))
	if manejadorFileError != nil {
		return respuestaApi, manejadorFileError
	}
	imagenOrigen, imagenOrigenError := os.Open(rutaOrigen)
	if imagenOrigenError != nil {
		return respuestaApi, imagenOrigenError
	}
	defer imagenOrigen.Close()
	_, copiaError := io.Copy(manejadorFile, imagenOrigen)
	if copiaError != nil {
		return respuestaApi, copiaError
	}

	manejadorForm.Close()

	peticion, peticionError := http.NewRequest("POST", "https://api.platerecognizer.com/v1/plate-reader/", form)
	if peticionError != nil {
		return respuestaApi, peticionError
	}
	peticion.Header.Set("Authorization", "Token "+token)
	peticion.Header.Set("Content-Type", manejadorForm.FormDataContentType())
	respuesta, respuestaError := cliente.Do(peticion)
	if respuestaError != nil {
		return respuestaApi, respuestaError
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode < 200 || respuesta.StatusCode > 299 {
		return respuestaApi, errors.New("status code incorrecto:" + respuesta.Status)
	}

	errorDecodificacion := json.NewDecoder(respuesta.Body).Decode(&respuestaApi)
	if errorDecodificacion != nil {
		return respuestaApi, errorDecodificacion
	}

	return respuestaApi, nil
}

func (m *Manejador) Procesar() {
	var progreso zenity.ProgressDialog
	var progresoError error
	if m.Configuracion.Dialogos {
		progreso, progresoError = zenity.Progress(
			zenity.Title("Procesando, por favor, espera"),
			zenity.Width(600),
			zenity.NoCancel())
		if progresoError != nil {
			zenity.Info(progresoError.Error(),
				zenity.Title("ERROR CRÍTICO"),
				zenity.Width(600),
				zenity.ErrorIcon)
			log.Fatalln(progresoError)
		}
	}

	for i, v := range m.Archivos {
		if m.Configuracion.Dialogos {
			porcentajeProgreso := (i + 1) * 100 / len(m.Archivos)
			mensajeProgreso := fmt.Sprintf("Procesando imagen: %v%v (%d de %d)", v.Nombre, v.Extension, i+1, len(m.Archivos))
			progreso.Text(mensajeProgreso)
			progreso.Value(porcentajeProgreso)
		}

		// Petición API
		rutaOrigen := v.Ruta
		nombre := fmt.Sprintf("%v%v", v.Nombre, v.Extension)
		rutaDestino := filepath.Join(m.Configuracion.RutaDestino, nombre)
		// Parece que la API solo reconoce jpg, transformar aquellas que no lo son
		if v.Extension != ".jpg" && v.Extension != ".jpeg" {
			nombreArchivoProvisional := fmt.Sprintf("%v.jpg", v.Nombre)
			rutaArchivoProvisional := filepath.Join(m.Configuracion.RutaDestino, nombreArchivoProvisional)
			comandoConvertir, comandoConvertirError := exec.Command(m.Configuracion.RutaImageMagick, "convert", rutaOrigen, rutaArchivoProvisional).CombinedOutput()
			if comandoConvertirError != nil {
				m.Errores = true
				informe := []string{nombre, string(comandoConvertir), rutaOrigen, "-", "-", "-"}
				m.Informe = append(m.Informe, informe)
				continue
			}
			rutaOrigen = rutaArchivoProvisional
			rutaDestino = rutaArchivoProvisional
		}

		resultados, resultadosError := peticionAPI(m.Cliente, rutaOrigen, rutaDestino, m.Configuracion.Token)
		if resultadosError != nil {
			m.Errores = true
			informe := []string{nombre, resultadosError.Error(), rutaOrigen, "-", "-", "-"}
			m.Informe = append(m.Informe, informe)
			continue
		}
		// Imagemagick
		var infoCoordenadas []string
		var totalResultados int
		origen := v.Ruta
		var procesamientoResultadosError string
		for _, resultado := range resultados.Results {
			totalResultados++
			infoC := fmt.Sprintf("[%d,%d,%d,%d]", resultado.Box.Xmin, resultado.Box.Ymin, resultado.Box.Xmax, resultado.Box.Ymax)
			infoCoordenadas = append(infoCoordenadas, infoC)

			argumentoRectangulo := fmt.Sprintf(`rectangle %d,%d %d,%d`, resultado.Box.Xmin, resultado.Box.Ymin, resultado.Box.Xmax, resultado.Box.Ymax)

			comando := exec.Command(m.Configuracion.RutaImageMagick, origen, "-draw", argumentoRectangulo, rutaDestino)
			origen = rutaDestino
			comandoError := comando.Run()
			if comandoError != nil {
				procesamientoResultadosError = comandoError.Error()
				break
			}

		}

		if procesamientoResultadosError != "" {
			m.Errores = true
			informe := []string{nombre, procesamientoResultadosError, rutaOrigen, "-", "-", "-"}
			m.Informe = append(m.Informe, informe)
			continue
		}

		// Si no ha habido resultados no se ha copiado la imagen
		if totalResultados == 0 {
			rutaDestino = ""
		}

		informe := []string{nombre, "", rutaOrigen, rutaDestino, fmt.Sprint(totalResultados), strings.Join(infoCoordenadas, "; ")}
		m.Informe = append(m.Informe, informe)
	}

	if m.Configuracion.Dialogos {
		progreso.Complete()
		progreso.Close()
	}
}
