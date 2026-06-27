package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	c.SetRequestTimeout(2 * time.Second)

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("✅ ÉXITO: %s (Status: %d)\n", r.Request.URL, r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("------------------------------------------------")
		fmt.Printf("❌ ERROR en: %s\n", r.Request.URL)
		fmt.Printf("   Status: %d | Error: %s\n", r.StatusCode, err)

		// 1. LEER EL CONTADOR DESDE LOS HEADERS
		// Los headers viajan con la petición, incluso si el contexto muere por timeout.
		intentosStr := r.Request.Headers.Get("X-Intentos")
		intentos := 0
		if intentosStr != "" {
			intentos, _ = strconv.Atoi(intentosStr)
		}

		// 2. VERIFICAR LÍMITE
		if (r.StatusCode >= 500 || r.StatusCode == 0) && intentos < 3 {
			intentos++
			fmt.Printf("   🔄 Reintentando... (Intento %d de 3)\n", intentos)

			// 3. GUARDAR EL NUEVO VALOR EN EL HEADER
			r.Request.Headers.Set("X-Intentos", strconv.Itoa(intentos))
			
			r.Request.Retry()
		} else {
			fmt.Println("   ⛔ Se acabaron los intentos o error fatal.")
		}
		fmt.Println("------------------------------------------------")
	})

	fmt.Println("--- Iniciando Prueba con Headers (Anti-Bucle) ---")
	// Probamos el caso difícil: Timeout
	c.Visit("https://httpbin.org/delay/5")
}
