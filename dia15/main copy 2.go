package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// ============================================================
	// CONCEPTO: Evaluate retorna JSON
	// Go no puede recibir un array de JS directamente.
	// Por eso: JS retorna JSON string → Go lo convierte a []string
	// ============================================================

	// Paso 1: declarar variable que recibirá el JSON como string
	var autoresJSON string
	var textosJSON  string
	var tagsJSON    string

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://quotes.toscrape.com/"),
		chromedp.WaitVisible(`.quote`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),

		// ✅ Evaluate — extraer TODOS los autores como array JSON
		// JSON.stringify convierte el array JS → string que Go puede leer
		chromedp.Evaluate(`
			JSON.stringify(
				Array.from(document.querySelectorAll('.author'))
				     .map(el => el.innerText.trim())
			)
		`, &autoresJSON),

		// ✅ Evaluate — extraer TODOS los textos
		chromedp.Evaluate(`
			JSON.stringify(
				Array.from(document.querySelectorAll('.text'))
				     .map(el => el.innerText.trim())
			)
		`, &textosJSON),

		// ✅ Evaluate — extraer TODOS los tags (primer tag de cada cita)
		chromedp.Evaluate(`
			JSON.stringify(
				Array.from(document.querySelectorAll('.quote'))
				     .map(quote => {
				         const tag = quote.querySelector('.tag');
				         return tag ? tag.innerText.trim() : 'sin-tag';
				     })
			)
		`, &tagsJSON),
	)

	if err != nil {
		log.Fatal("❌ Error:", err)
	}

	// ============================================================
	// Paso 2: convertir JSON string → slice de Go []string
	// json.Unmarshal "desempaca" el JSON en una variable Go
	// ============================================================
	var autores []string
	var textos  []string
	var tags    []string

	json.Unmarshal([]byte(autoresJSON), &autores)
	json.Unmarshal([]byte(textosJSON),  &textos)
	json.Unmarshal([]byte(tagsJSON),    &tags)

	// ============================================================
	// Imprimir resultados
	// ============================================================
	fmt.Printf("\n📊 Total citas encontradas: %d\n", len(autores))
	fmt.Println("═══════════════════════════════════════════════")

	for i := range autores {
		fmt.Printf("\n[%d] ✍️  %s\n", i+1, autores[i])
		fmt.Printf("    🏷️  Tag: %s\n", tags[i])

		// Mostrar solo los primeros 60 chars del texto
		texto := textos[i]
		if len(texto) > 60 {
			texto = texto[:60] + "..."
		}
		fmt.Printf("    💬 %s\n", texto)
	}

	fmt.Println("\n✅ Ejercicio 2 completado!")
}
