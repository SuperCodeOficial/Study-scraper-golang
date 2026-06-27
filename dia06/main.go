package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings" // 1. IMPORTANTE: Necesitamos esto para limpiar el texto

	"github.com/gocolly/colly/v2"
)

func main() {
	nombreCarpeta := "imagenes_descargadas"
	os.Mkdir(nombreCarpeta, os.ModePerm)

	c := colly.NewCollector()

	c.OnHTML("div.thumbnail", func(e *colly.HTMLElement) {
		imgRelativa := e.ChildAttr("img", "src")
		imgUrl := e.Request.AbsoluteURL(imgRelativa)

		// --- CORRECCIÓN AQUÍ ---
		
		// 1. Obtenemos el nombre crudo: "Asus VivoBook"
		nombreCrudo := e.ChildText("a.title")

		// 2. Limpiamos el nombre
		// Reemplazamos espacios por guiones para que se vea prolijo: "Asus-VivoBook"
		nombreLimpio := strings.ReplaceAll(nombreCrudo, " ", "-")
		// Opcional: Si hubiera barras /, habría que quitarlas también para no romper la ruta
		nombreLimpio = strings.ReplaceAll(nombreLimpio, "/", "-")

		// 3. Agregamos la extensión manualmente
		// Asumimos que son jpg. (En un nivel avanzado detectaríamos si es png o jpg automáticamente)
		nombreArchivo := nombreLimpio + ".jpg"

		rutaCompleta := nombreCarpeta + "/" + nombreArchivo

		fmt.Println("--> Descargando:", nombreArchivo)

		err := descargarArchivo(imgUrl, rutaCompleta)
		if err != nil {
			fmt.Println("Error descargando:", err)
		}
	})

	c.Visit("https://webscraper.io/test-sites/e-commerce/allinone/computers/laptops")
}

func descargarArchivo(url string, rutaDestino string) error {
	respuesta, err := http.Get(url)
	if err != nil {
		return err
	}
	defer respuesta.Body.Close()

	archivo, err := os.Create(rutaDestino)
	if err != nil {
		return err
	}
	defer archivo.Close()

	_, err = io.Copy(archivo, respuesta.Body)
	return err
}
