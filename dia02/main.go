package main

import (
	"fmt"

	"github.com/gocolly/colly/v2" // Importamos Colly
)

func main() {
	// 1. Crear el "Recolector" (Collector)
	// Es el robot que hará el trabajo por nosotros.
	c := colly.NewCollector(
		// Aquí podríamos configurar cosas como User-Agent, Proxies, etc.
		// Por ahora, lo dejamos por defecto.
	)

	// 2. Definir los "Callbacks" (Eventos)
	// Colly funciona por eventos: "Cuando pase X, haz Y".

	// Evento: OnRequest -> Antes de hacer la petición
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitando:", r.URL)
	})

	// Evento: OnError -> Si algo sale mal
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Algo salió mal:", err)
	})

	// Evento: OnResponse -> Cuando recibimos la respuesta (HTML completo)
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Página recibida. Status:", r.StatusCode)
	})

	// Evento: OnHTML -> Cuando encuentra un elemento específico
	// Aquí le decimos: "Busca la etiqueta <title> y dame su texto"
	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("-----------------------------------")
		fmt.Println("TÍTULO ENCONTRADO:", e.Text)
		fmt.Println("-----------------------------------")
	})

	// 3. ¡Lanzar el robot!
	// Todo lo anterior fue configuración. Aquí empieza la acción.
	fmt.Println("🤖 Iniciando el robot...")
	c.Visit("https://webscraper.io/test-sites/e-commerce/allinone")
}
