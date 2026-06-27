package main

import (
	"encoding/json"
	"fmt"

	"github.com/gocolly/colly/v2"
)

// 1. EL PLANO (STRUCTS)
// Definimos cómo se ve un solo producto en el JSON
type Product struct {
	Title       string  `json:"title"`       // Mapea el campo "title" del JSON a esta variable
	Price       float64 `json:"price"`       // Mapea "price"
	Description string  `json:"description"`
	Category string `json:"category"` 
	Stock int64 `json:"stock"`// Mapea "description"
}

// Definimos la respuesta completa (que contiene una lista de productos)
type APIResponse struct {
	Total    int       `json:"total"`
	Products []Product `json:"products"` // Una lista de la estructura Product
}

func main() {
	c := colly.NewCollector()

	// 2. YA NO USAMOS OnHTML
	// Como la respuesta no es HTML, sino texto plano (JSON), usamos OnResponse
	c.OnResponse(func(r *colly.Response) {
		// r.Body contiene los bytes del JSON crudo

		// Preparamos la variable donde guardaremos los datos
		var respuesta APIResponse

		// 3. DECODIFICAR (UNMARSHAL)
		// Convertimos los bytes raros a nuestra estructura de Go
		err := json.Unmarshal(r.Body, &respuesta)
		if err != nil {
			fmt.Println("Error al leer JSON:", err)
			return
		}

		fmt.Printf("✅ Se encontraron %d productos en la API:\n", respuesta.Total)
		fmt.Println("------------------------------------------------")

		// 4. RECORRER LOS DATOS LIMPIOS
		for _, p := range respuesta.Products {
			fmt.Printf("📦 %s\n", p.Title)
			fmt.Printf("   💰 Precio: $%.2f\n", p.Price)
			fmt.Printf("   💰 Categoria: %s\n", p.Category)
			fmt.Printf("   💰 Cantidad Stock: %d\n", p.Stock)
			fmt.Printf("   📝 Info: %s\n", p.Description[:50]+"...") // Solo los primeros 50 caracteres
			fmt.Println("-")
		}
	})

	fmt.Println("--- Interceptando API JSON ---")
	// Visitamos directamente la URL de la API, no la web visual
	c.Visit("https://dummyjson.com/products?limit=5")
}
