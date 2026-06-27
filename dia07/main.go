package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {

	// 1. Configurar CSV
	nombreArchivo := "productos1.csv"
	archivo, err := os.Create(nombreArchivo)
	if err != nil {
		log.Fatalf("No se pudo crear el archivo, error: %q", err)
		return
	}
	defer archivo.Close()

	escritor := csv.NewWriter(archivo)
	defer escritor.Flush() 

	escritor.Write([]string{"Nombre", "Precio", "Descripcion", "Reviews", "Rating"})

    // 2. Configurar Carpeta Imágenes
	nombreCarpeta := "imagenes_descargadas"
	os.Mkdir(nombreCarpeta, os.ModePerm)

    // 3. Configurar Collector
	c := colly.NewCollector()

		c.OnHTML("div.thumbnail", func(e *colly.HTMLElement) {
		imgRelativa := e.ChildAttr("img", "src")
		imgUrl := e.Request.AbsoluteURL(imgRelativa)

		nombreCrudo := e.ChildText("a.title")
		nombreLimpio := strings.ReplaceAll(nombreCrudo, " ", "-")
		nombreLimpio = strings.ReplaceAll(nombreLimpio, "/", "-")
		nombreArchivo := nombreLimpio + ".jpg"
		rutaCompleta := nombreCarpeta + "/" + nombreArchivo

		nombre := e.ChildText("a.title")
		precio := e.ChildText("h4.pull-right.price")
		descripcion := e.ChildText("p.description")
		reviews := e.ChildText(".review-count")
		rating := e.ChildAttr("p[data-rating]", "data-rating")
		escritor.Write([]string{nombre, precio, descripcion, reviews, rating})

		fmt.Println("--> Descargando:", nombreArchivo)
		err := descargarArchivo(imgUrl, rutaCompleta)
		if err != nil {
			fmt.Println("Error descargando:", err)
		}
		fmt.Println("Guardado:", nombre)
	})

    // Callback para PAGINACIÓN (Encontrar el botón "Next" o los números)
c.OnHTML("a[rel='next']", func(e *colly.HTMLElement) {
    link := e.Attr("href")
    fmt.Println("--> Yendo a la siguiente página:", link)
    e.Request.Visit(link)
})

		c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitando URL:", r.URL)
	})
    
	c.Visit("https://webscraper.io/test-sites/e-commerce/static/computers/tablets")
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
