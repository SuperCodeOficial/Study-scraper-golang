package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	// IMPORTANTE: Para convertir texto a número
	// Para limpiar el texto ($)
	"github.com/gocolly/colly/v2"
)

func main() {
	totalDescargados := 0
	totalSaltados := 0
	nombreArchivo := "laptops_filtradas1.csv"
	archivo, err := os.Create(nombreArchivo)
	if err != nil {
		log.Fatalf("No fue posible crear el archivo, error: %q", err)
		return
	}

	defer archivo.Close()

	escritor := csv.NewWriter(archivo)
	defer escritor.Flush()

	escritor.Write([]string{"Nombre", "Precio", "Marca"})

	nombreCarpeta := "backup_imagenes"
	os.Mkdir(nombreCarpeta, os.ModePerm)

	c := colly.NewCollector()

	c.OnHTML("div.thumbnail", func(e *colly.HTMLElement) {
		imgRelativa := e.ChildAttr("img", "src")
		imgUrl := e.Request.AbsoluteURL(imgRelativa)

		nombreCrudo := e.ChildText("a.title")
		nombreLimpio := strings.ReplaceAll(nombreCrudo, " ", "-")
		nombreLimpio = strings.ReplaceAll(nombreLimpio, "/", "-")
		nombreArchivo := nombreLimpio + ".jpg"
		rutaCompleta := nombreCarpeta + "/" + nombreArchivo
		
		if archivoExiste(rutaCompleta) {
			fmt.Println("Saltando: (ya existe)", nombreArchivo)
			totalSaltados++
		} else {
			fmt.Println("--> Descargando Nuevo:", nombreArchivo)
			err := descargarArchivo(imgUrl, rutaCompleta)
			if err != nil {
				fmt.Println("Error descargando:", err)
			}
			totalDescargados++
		}
		nombre := e.ChildText("a.title")
		precioTexto := e.ChildText("h4.price") 
		
		precioLimpio := strings.ReplaceAll(precioTexto, "$", "")
		precioFinal, err := strconv.ParseFloat(precioLimpio, 64)
		
		if err != nil {
			fmt.Println("Error convirtiendo el precio:", err)
			return
		}
		fmt.Println("El precio de "+nombre, "es igual a: ",+precioFinal)
		
		if precioFinal < 100.0 {
			marca := ""
			if strings.Contains(strings.ToLower(nombre), "lenovo") {
			marca = "Lenovo"
			
			} else {
				marca = "Otra"
			}
			escritor.Write([]string{nombre, fmt.Sprintf("%.2f", precioFinal), marca})

			
		}
		
		
	})

	c.Visit("https://webscraper.io/test-sites/e-commerce/allinone/computers/laptops")

	fmt.Println("=======================================")
	fmt.Println("RESUMEN DEL PROCESO")
	fmt.Printf("✅ Imágenes nuevas descargadas: %d\n", totalDescargados)
	fmt.Printf("♻️ Imágenes ya existentes:      %d\n", totalSaltados)
	fmt.Printf("📊 Total procesado:             %d\n", totalDescargados+totalSaltados)
	fmt.Println("=======================================")
}
func archivoExiste(ruta string) bool {
	_, err := os.Stat(ruta)
	if os.IsNotExist(err) {
		return false
	}
	return true
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