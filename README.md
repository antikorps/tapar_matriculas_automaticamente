# tapar_matriculas_automaticamente
Utilidad para tapar matrículas de coches automáticamente. 

La identificación de la matrícula se hace a través de la web [Plate Recognizer](https://platerecognizer.com/), por lo que será necesario contar con un token válido de su web. Ofrece un plan gratuito de hasta 2500 matrículas por mes en el que no es necesario incorporar un número de tarjeta.

Tras el procesamiento se ofrece un informe en formato CSV en el que se resume todo el proceso de ejecución, recopilando las coordenadas de las matrículas que ofrece Plate Recognizer en su respuesta.

## Instalación y uso
Para su ejecución simplemente es necesario descargar el ZIP del sistema operativo en el que se va a ejecutar de la carpeta bin. Al descomprimir este archivo aparecerá la siguiente estructura de archivos:
- configuracion: es un directorio con el archivo de configuración y el ejecutable de imagemagick necesario para tapar las matrículas
- tapar_matriculas_automatica: es el ejecutable del programa.

Antes de su ejecución es necesario incorporar un token válido para Plate Recognizer. El token que aparece por defecto es el que he utilizado para hacer pruebas, al momento de usarse quizá ya no sea operativo. Una vez creada la cuenta se puede obtener desde el siguiente [enlace](https://app.platerecognizer.com/products/snapshot-cloud/)

![Token](https://i.imgur.com/I1JQUyU.png)

Si eres usuario de GNU/Linux recuerda que quizá debas otorgar permisos de ejecución a los ejecutables:
```bash
sudo cd /home/TU_USUARIO/RUTA_HASTA_EL_ZIP_DESCOMPRIMIDO && chmod +x tapar_matriculas_automaticamente && chmod +x configuracion/imagemagick
```

### Uso por línea de comandos
Para usarse por línea de comandos se debe pasar -origen (ruta completa a la imagen o al directorio que contiene las imágenes) y -destino (ruta a un directorio en el que se guardarán las nuevas imágenes creadas).
```bash
 ./tapar_matriculas_automaticamente -origen /home/TU_USUARIO/RUTA_MATRICULAS/ -destino /home/TU_USUARIO/RUTA_MATRICULAS_TAPADAS
```

### Uso por interfaz de diálogos
Si se ejecuta directamente el binario sin ningún tipo de argumento se iniciará la interfaz de diálogos.

![Presentación](https://i.imgur.com/ojFU9Uy.png)

Después de la presentación, será necesario escoger el tipo de procesamiento que se desea:
- individual: permite seleccionar una imagen y procesar únicamente ese archivo.
- multiple: permite seleccionar un directorio y buscará todas las imágenes en esa carpeta y subcarpetas.

![Tipo de procesamiento](https://i.imgur.com/bGF0V9q.png)

A continuación deberá seleccionarse el origen y el destino (carpeta en la que se guardarán las imágenes convertidas)

Posteriormente apareceŕa un diálogo de proceso, debe esperarse a que se complete sin pulsar "Aceptar"
![Progreso](https://i.imgur.com/2xu1jpy.png)

Finalmente aparecerá una ventana de información con el resultado de la operación y el recordatorio del informe generado.

![Informe](https://i.imgur.com/UT0BleQ.png)

## Ejemplos
![Ejemplo](https://i.imgur.com/IP00U4e.png)

![Ejemplo](https://i.imgur.com/hXgHh4i.png)

![Ejemplo](https://i.imgur.com/SI0zLav.png)

![Ejemplo](https://i.imgur.com/LoPTtVx.png)

![Ejemplo](https://i.imgur.com/w15ODJB.png)

![Ejemplo](https://i.imgur.com/JMaYrWx.png)

![Ejemplo](https://i.imgur.com/AfdmWsU.png)