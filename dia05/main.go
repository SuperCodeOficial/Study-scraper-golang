package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()

	// --- ESTRATEGIA DE NAVEGACIÓN ---

	// 1. DETECTAR ENLACES (En la lista de productos)
	// Buscamos los enlaces que tienen la clase "title"
	c.OnHTML("a.title", func(e *colly.HTMLElement) {
		// Extraemos el atributo "href" (la dirección a donde lleva el link)
		link := e.Attr("href")
		
		fmt.Println("--> Encontré un enlace a:", link)

		// LA MAGIA: Le decimos al robot "Visita ese enlace"
		// Visit toma el link relativo (/test-sites/...) y lo combina con el dominio
		e.Request.Visit(link)
	})

	// 2. EXTRAER DATOS (En la página de detalle)
	// Este código se ejecutará cuando el robot entre a la página del producto.
	// En la página de detalle, el precio no está en un "div.caption", 
	// sino que suele estar más grande. Buscamos detalles específicos.
	c.OnHTML("div.caption", func(e *colly.HTMLElement) {
		// Nota: En la página interna, la estructura es similar, 
		// así que podemos reusar lógica o ser más específicos.
		// Para este ejemplo, imprimiremos el nombre para confirmar que entramos.
		
		nombre := e.ChildText("h4:nth-of-type(2)") // El segundo h4 es el nombre en la vista detalle
		precio := e.ChildText("h4.pull-right.price")

		if nombre != "" {
			fmt.Println("   [DENTRO DEL PRODUCTO] Nombre:", nombre, "| Precio:", precio)
		}
	})

	// Log para ver qué URL está cargando realmente
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitando URL:", r.URL)
	})

	// Empezamos en la lista de Laptops
	c.Visit("https://webscraper.io/test-sites/e-commerce/allinone/computers/laptops")
}
