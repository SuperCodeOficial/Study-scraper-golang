package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

// 1. DEFINIMOS LA ESTRUCTURA (La "Ficha" del dato)
// Esto le dice a Go qué forma tiene un producto.
type Producto struct {
	Nombre string
	Precio string
	Descripcion string
}

func main() {
	c := colly.NewCollector()

	// 2. APUNTAMOS AL CONTENEDOR
	// En lugar de buscar solo el titulo, buscamos la "caja" que envuelve a CADA producto.
	// En esta web, cada producto está dentro de un div con clase "caption"
	c.OnHTML("div.caption", func(e *colly.HTMLElement) {
		// 3. CREAMOS UNA NUEVA FICHA VACÍA
		p := Producto{}

		// 4. LLENAMOS LA FICHA (Buscamos DENTRO de la caja)
		// Usamos ChildText para buscar hijos de este elemento específico
		p.Nombre = e.ChildText("a.title")
		p.Precio = e.ChildText("h4.pull-right.price")
		p.Descripcion = e.ChildText("p.description")

		// 5. IMPRIMIMOS LA FICHA COMPLETA
		fmt.Println("-----------------------------")
		fmt.Println("Nombre:", p.Nombre)
		fmt.Println("Precio:", p.Precio)
		fmt.Println("Desc  :", p.Descripcion)
		fmt.Println("-----------------------------")
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitando:", r.URL)
	})

	c.Visit("https://webscraper.io/test-sites/e-commerce/allinone/phones")
}
