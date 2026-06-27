package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Iniciamos el cronómetro para ver cuánto tarda
	inicio := time.Now()

	c := colly.NewCollector(
		// 1. ACTIVAR MODO ASÍNCRONO
		// Esto le dice a Colly: "No esperes a terminar una visita para empezar la siguiente"
		colly.Async(true),
	)

	// 2. REGLAS DE LÍMITE (¡OBLIGATORIO EN ASYNC!)
	// Si no pones esto, Colly lanzará 1000 peticiones en 1 milisegundo y te banearán la IP.
	// Aquí controlamos el caos.
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*", // Aplica a cualquier dominio
		Parallelism: 4,   // Número de "repartidores" (hilos) simultáneos
		Delay:       1 * time.Second, // Espera entre peticiones (cortesía)
	})

	// Callback: Cuando encontramos un producto
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		nombre := e.ChildText("h2")
		// Imprimimos algo corto para no saturar la consola
		fmt.Printf("   📦 Encontrado: %s\n", nombre)
	})

	// Callback: Para saber qué página estamos visitando
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("🚀 Visitando:", r.URL)
	})

	// 3. EL BUCLE DE URLs
	// Vamos a encolar 5 páginas.
	// En modo normal, tardaría unos 10-15 segundos.
	// En modo Async, tardará mucho menos.
	for i := 1; i <= 10; i++ {
		url := fmt.Sprintf("https://scrapeme.live/shop/page/%d/", i)
		c.Visit(url)
	}


	// 4. ESPERAR A LOS ROBOTS (WAIT)
	// Esto es CRUCIAL. Como el código es asíncrono, el programa principal (main)
	// podría terminar y cerrarse antes de que los robots vuelvan con los datos.
	// c.Wait() dice: "No te cierres hasta que todos los hilos terminen".
	c.Wait()

	// Calculamos el tiempo total
	tiempoTotal := time.Since(inicio)
	fmt.Printf("\n⏱️ Tiempo total de ejecución: %s\n", tiempoTotal)
}
