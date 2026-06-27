package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gocolly/colly/v2"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Mobile/15E148 Safari/604.1",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	c := colly.NewCollector()
	c.AllowURLRevisit = true

	// 1. PREPARAR LA PETICIÓN
	c.OnRequest(func(r *colly.Request) {
		randomIndex := rand.Intn(len(userAgents))
		ua := userAgents[randomIndex]

		r.Headers.Set("User-Agent", ua)
		// r.Headers.Set("X-Mi-Robot", "SuperMaxBot-v1") // <--- Tu header personalizado

		fmt.Printf("🎭 Enviando como: %s...\n", ua[0:30]) // Imprimimos solo el inicio para no ensuciar
	})

	// 2. RECIBIR LA RESPUESTA (RAW)
	// Usamos OnResponse en vez de OnHTML para asegurarnos de ver TODO
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("✅ Respuesta recibida (Status 200)")
		
		// Convertimos los bytes de la respuesta a String para leerla
		cuerpo := string(r.Body)
		fmt.Println(cuerpo)
		fmt.Println("------------------------------------------------")
	})

	// 3. CAPTURAR ERRORES
	// Si algo falla, queremos saberlo
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("❌ Error:", err)
	})

	url := "https://httpbin.org/headers"

	fmt.Println("--- Intento 1 ---")
	c.Visit(url)
	
	fmt.Println("--- Intento 2 ---")
	c.Visit(url)

	fmt.Println("--- Intento 3 ---")
	c.Visit(url)
}
