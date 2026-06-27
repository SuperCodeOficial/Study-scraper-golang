package main

import (
	"encoding/csv" // 1. Importamos la herramienta para CSV
	"fmt"
	"log"
	"os" // 2. Importamos herramientas para manejar archivos del sistema (Operating System)

	"github.com/gocolly/colly/v2"
)

func main() {
	// --- PARTE A: PREPARAR EL ARCHIVO ---
	
	// 3. Creamos el archivo "productos.csv"
	nombreArchivo := "productos.csv"
	archivo, err := os.Create(nombreArchivo)
	if err != nil {
		log.Fatalf("No se pudo crear el archivo, error: %q", err)
		return
	}
	// "defer" significa: "Ejecuta esto justo antes de que termine la función main"
	// Es importante cerrar el archivo al final para guardar los cambios.
	defer archivo.Close()

	// 4. Creamos el "Escritor" (el bolígrafo)
	escritor := csv.NewWriter(archivo)
	defer escritor.Flush() // Asegura que se escriban todos los datos pendientes al final

	// 5. Escribimos el ENCABEZADO (La primera fila del Excel)
	escritor.Write([]string{"Nombre", "Precio", "Descripción"})

	// --- PARTE B: CONFIGURAR EL ROBOT (Igual que antes) ---

	c := colly.NewCollector()

	c.OnHTML("div.caption", func(e *colly.HTMLElement) {
		// Extraemos los datos
		nombre := e.ChildText("a.title")
		precio := e.ChildText("h4.pull-right.price")
		descripcion := e.ChildText("p.description")

		// 6. GUARDAMOS EN EL ARCHIVO
		// En lugar de fmt.Println, usamos escritor.Write
		// Le pasamos una lista de strings []string{...}
		escritor.Write([]string{nombre, precio, descripcion})
		
		// Imprimimos un mensajito solo para saber que está trabajando
		fmt.Println("Guardado:", nombre)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitando:", r.URL)
	})

	// Usamos la URL de teléfonos que sabemos que funciona
	c.Visit("https://webscraper.io/test-sites/e-commerce/allinone/phones")
}
