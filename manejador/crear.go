package manejador

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ncruces/zenity"
)

func (c *Configuracion) configurarRutas() {
	rutaEjecutable, rutaEjecutableError := os.Executable()
	if rutaEjecutableError != nil {
		mensajeError := fmt.Sprintf("no se ha podido obtener la ruta del ejecutable: %v", rutaEjecutableError)
		if c.Dialogos {
			zenity.Info(mensajeError,
				zenity.Title("ERROR CRÍTICO"),
				zenity.Width(600),
				zenity.ErrorIcon)
		}
		log.Fatalln(mensajeError)
	}
	rutaRaiz := filepath.Dir(rutaEjecutable)
	rutaCarpetaConfiguracion := filepath.Join(rutaRaiz, "configuracion")

	c.RutaConfiguracion = filepath.Join(rutaCarpetaConfiguracion, "configuracion.txt")

	rutaImageMagick := filepath.Join(rutaCarpetaConfiguracion, "imagemagick")
	if runtime.GOOS == "windows" {
		rutaImageMagick = filepath.Join(rutaCarpetaConfiguracion, "imagemagick.exe")
	}
	c.RutaImageMagick = rutaImageMagick
}

func (c *Configuracion) leerConfiguracion() {
	archivo, archivoError := os.Open(c.RutaConfiguracion)
	if archivoError != nil {
		mensajeError := fmt.Sprintf("no se ha podido leer el archivo de configuración: %v", archivoError)
		if c.Dialogos {
			zenity.Info(mensajeError,
				zenity.Title("ERROR CRÍTICO"),
				zenity.Width(600),
				zenity.ErrorIcon)
		}
		log.Fatalln(mensajeError)
	}
	defer archivo.Close()
	escaner := bufio.NewScanner(archivo)
	escaner.Split(bufio.ScanLines)
	var indice int
	for escaner.Scan() {
		indice++
		linea := escaner.Text()
		if indice == 1 {
			if !strings.HasPrefix(linea, "TOKEN=") {
				mensajeError := "no se ha encontrado la clave TOKEN= en la primera línea del archivo de configuración"
				if c.Dialogos {
					zenity.Info(mensajeError,
						zenity.Title("ERROR CRÍTICO"),
						zenity.Width(600),
						zenity.ErrorIcon)
				}
				log.Fatalln(mensajeError)
			}
			token := strings.TrimPrefix(linea, "TOKEN=")
			if token == "" {
				mensajeError := "no se ha encontrado ningún valor para la clave TOKEN en la primera línea del archivo de configuración"
				if c.Dialogos {
					zenity.Info(mensajeError,
						zenity.Title("ERROR CRÍTICO"),
						zenity.Width(600),
						zenity.ErrorIcon)
				}
				log.Fatalln(mensajeError)
			}
			c.Token = token
		}

		if indice == 2 {
			if !strings.HasPrefix(linea, "EXTENSIONES=") {
				mensajeError := "no se ha encontrado la clave EXTENSIONES= en la tercera línea del archivo de configuración"
				if c.Dialogos {
					zenity.Info(mensajeError,
						zenity.Title("ERROR CRÍTICO"),
						zenity.Width(600),
						zenity.ErrorIcon)
				}
				log.Fatalln(mensajeError)
			}
			extensiones := strings.TrimPrefix(linea, "EXTENSIONES=")
			if extensiones == "" {
				mensajeError := "no se ha encontrado ningún valor para la clave EXTENSIONES en la tercera línea del archivo de configuración"
				if c.Dialogos {
					zenity.Info(mensajeError,
						zenity.Title("ERROR CRÍTICO"),
						zenity.Width(600),
						zenity.ErrorIcon)
				}
				log.Fatalln(mensajeError)
			}
			var extensionesValidas []string
			ext := strings.Split(extensiones, ",")
			for _, v := range ext {
				if !strings.HasPrefix(v, ".") {
					v = "." + v
				}
				extension := strings.ToLower(v)
				extension = strings.TrimSpace(extension)
				extensionesValidas = append(extensionesValidas, extension)
			}
			c.ExtensionesValidas = extensionesValidas
		}
	}
}

func (c *Configuracion) recuperarRutas() {
	if c.Dialogos {
		var origen, destino string
		var origenError, destinoError error

		presentacionMensaje := `Ten en cuenta las siguientes cuestiones:

	- Respecto al reconocimiento de las matrículas:
	   - Depende de platerecognizer.com
	   - Es necesario disponer de token
	   - Existe un plan gratuito de 2500 imágenes por mes

	- Respecto al funcionamiento del programa:
	   - Debe cumplimentarse la linea TOKEN del archivo configuracion.txt
	   - Se puede ver en https://app.platerecognizer.com/products/snapshot-cloud/
	   - Modifica el resto de líneas solamente si sabes que implica

	- Selecciona el tipo de procesamiento:
	   - individual (una sola imgen)
	   - múltiple (todas las imágenes de una carpeta y sus subdirectorios)

	- Selecciona la carpeta en la que guardar las nuevas imágenes`

		presentacionError := zenity.Info(presentacionMensaje,
			zenity.Title("Presentación"),
			zenity.Width(600),
			zenity.InfoIcon)
		if presentacionError != nil {
			if presentacionError == zenity.ErrCanceled {
				os.Exit(0)
			}
			log.Fatalln("error crítico en el diálogo de presentación", presentacionError)
		}

		tipo, tipoError := zenity.List(
			"Selecciona el tipo de procesamiento:",
			[]string{"individual", "múltiple"},
			zenity.Title("Tipo de extracción"),
			zenity.DisallowEmpty(),
		)
		if tipoError != nil {
			log.Fatalln("error crítico en el diálogo de tipo")
		}

		if tipo == "individual" {
			var extensionesSelector []string
			for _, ext := range c.ExtensionesValidas {
				extensionSelector := fmt.Sprintf("*%v", ext)
				extensionesSelector = append(extensionesSelector, extensionSelector)
			}
			origen, origenError = zenity.SelectFile(
				zenity.FileFilters{
					{Name: "Imágenes válidas", Patterns: extensionesSelector, CaseFold: true},
				})
			if origenError != nil {
				if origenError == zenity.ErrCanceled {
					os.Exit(0)
				}
				log.Fatalln("error en el diálogo de selección individual", origenError)
			}
		}
		if tipo == "múltiple" {
			origen, origenError = zenity.SelectFile(
				zenity.Directory())
			if origenError != nil {
				if origenError == zenity.ErrCanceled {
					os.Exit(0)
				}
				log.Fatalln("error en el diálogo de selección múltiple", origenError)
			}
		}

		destino, destinoError = zenity.SelectFile(
			zenity.Directory())
		if destinoError != nil {
			if destinoError == zenity.ErrCanceled {
				os.Exit(0)
			}
			log.Fatalln("error en el diálogo de destino", destinoError)
		}

		c.RutaOrigen = origen
		c.RutaDestino = destino
	}
}

func (c *Configuracion) RecuperarArchivos() []Archivo {
	var archivos []Archivo

	destino, destinoError := os.Stat(c.RutaDestino)
	if destinoError != nil {
		mensajeError := fmt.Sprintf("no se ha podido analizar la ruta incorporada para el destino: %v", destinoError)
		if c.Dialogos {
			zenity.Info(mensajeError,
				zenity.Title("ERROR CRÍTICO"),
				zenity.Width(600),
				zenity.ErrorIcon)
		}
		log.Fatalln(mensajeError)
	}
	if !destino.IsDir() {
		mensajeError := "la ruta de destino debe ser un directorio o carpeta"
		if c.Dialogos {
			zenity.Info(mensajeError,
				zenity.Title("ERROR CRÍTICO"),
				zenity.Width(600),
				zenity.ErrorIcon)
		}
		log.Fatalln(mensajeError)
	}

	origen, origenError := os.Stat(c.RutaOrigen)
	if origenError != nil {
		mensajeError := fmt.Sprintf("no se ha podido analizar la ruta incorporada en el origen: %v", origenError)
		if c.Dialogos {
			zenity.Info(mensajeError,
				zenity.Title("ERROR CRÍTICO"),
				zenity.Width(600),
				zenity.ErrorIcon)
		}
		log.Fatalln(mensajeError)
	}
	if !origen.IsDir() {
		// Procesamiento individual
		archivos = append(archivos, Archivo{
			Extension: strings.ToLower(filepath.Ext(c.RutaOrigen)),
			Nombre:    nombreDeArchivo(c.RutaOrigen),
			Ruta:      c.RutaOrigen,
		})
	} else {
		// Procesamiento múltiple
		directorioError := filepath.Walk(c.RutaOrigen, func(ruta string, info os.FileInfo, err error) error {
			if err != nil {
				mensajeError := fmt.Sprintf("error crítico recorriendo el directorio proporcionado: %v", err)
				return errors.New(mensajeError)
			}
			if info.IsDir() {
				return nil
			}
			extension := strings.ToLower(filepath.Ext(ruta))
			if !contiene(c.ExtensionesValidas, extension) {
				return nil
			}
			archivos = append(archivos, Archivo{
				Extension: strings.ToLower(filepath.Ext(ruta)),
				Nombre:    nombreDeArchivo(ruta),
				Ruta:      ruta,
			})
			return nil
		})
		if directorioError != nil {
			if c.Dialogos {
				zenity.Info(directorioError.Error(),
					zenity.Title("ERROR CRÍTICO"),
					zenity.Width(600),
					zenity.ErrorIcon)
			}
			log.Fatalln(directorioError)
		}
	}

	return archivos

}

func Crear(origen, destino string) Manejador {
	var dialogos bool
	if origen == "" || destino == "" {
		dialogos = true
	}

	var configuracion Configuracion
	configuracion.RutaOrigen = origen
	configuracion.RutaDestino = destino
	configuracion.Dialogos = dialogos

	configuracion.configurarRutas()
	configuracion.leerConfiguracion()
	configuracion.recuperarRutas()

	archivos := configuracion.RecuperarArchivos()

	var informe [][]string

	informe = append(informe, []string{"imagen", "error", "ruta_origen", "ruta_destino", "matriculas_encontradas", "matriculas_coordenadas"})

	return Manejador{
		Archivos: archivos,
		Cliente: http.Client{
			Timeout: 16 * time.Second,
		},
		Configuracion: configuracion,
		Informe:       informe,
	}
}
