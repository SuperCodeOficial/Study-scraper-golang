package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()

	// --- AQUÍ ESTÁ LA MAGIA DE HOY ---
	// Definimos una regla de límites (LimitRule)
	c.Limit(&colly.LimitRule{
		// 1. DomainGlob: ¿A qué webs aplica esta regla?
		// "*" significa "a todas las webs".
		// Podrías poner "en.wikipedia.org" para limitar solo esa.
		DomainGlob: "*",

		// 2. Delay: Tiempo MÍNIMO de espera entre peticiones
		Delay: 2 * time.Second,

		// 3. RandomDelay: Tiempo EXTRA aleatorio
		// La espera total será: Delay + (un número entre 0 y RandomDelay)
		// En este caso: entre 2 y 5 segundos (2s base + hasta 3s extra)
		RandomDelay: 3 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		// Imprimimos la hora exacta para que veas la pausa
		fmt.Printf("⏰ %s - Visitando: %s\n", time.Now().Format("15:04:05"), r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("   ✅ Recibido!")
	})

	// Vamos a visitar 5 páginas seguidas
	// Usamos httpbin.org/anything/1, /2, etc. para simular productos distintos
	fmt.Println("--- Iniciando Paseo Lento ---")
	for i := 1; i <= 5; i++ {
		url := fmt.Sprintf("https://httpbin.org/anything/%d", i)
		c.Visit(url)
	}
	fmt.Println("--- Fin del Paseo ---")
}
