package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type AutorDetalle struct {
	Nombre      string
	FechaNac    string
	LugarNac    string
	Descripcion string
	URLPerfil   string
}

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

	var autor AutorDetalle

	err := chromedp.Run(ctx,

		// ── PÁGINA 1: Lista de citas ──────────────────────────
		chromedp.Navigate("https://quotes.toscrape.com/"),
		chromedp.WaitVisible(`.quote`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),

		// ✅ CLICK — hacer clic en el nombre del primer autor
		// Esto hace que Chrome navegue a su página de perfil
		chromedp.Click(`.quote:nth-child(1) .author + a`, chromedp.ByQuery),

		// ── PÁGINA 2: Perfil del autor ────────────────────────
		// Esperar que cargue la nueva página
		chromedp.WaitVisible(`.author-title`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),

		// Capturar URL de la nueva página
		chromedp.Location(&autor.URLPerfil),

		// Extraer datos del perfil
		chromedp.Text(`.author-title`, &autor.Nombre, chromedp.ByQuery),
		chromedp.Text(`.author-born-date`, &autor.FechaNac, chromedp.ByQuery),
		chromedp.Text(`.author-born-location`, &autor.LugarNac, chromedp.ByQuery),
		chromedp.Text(`.author-description`, &autor.Descripcion, chromedp.ByQuery),
	)

	if err != nil {
		log.Fatal("❌ Error:", err)
	}

	// Recortar descripción si es muy larga
	desc := autor.Descripcion
	if len(desc) > 200 {
		desc = desc[:200] + "..."
	}

	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║     PERFIL DEL AUTOR — Ejercicio 3   ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Printf("\n👤 Nombre:       %s\n", autor.Nombre)
	fmt.Printf("📅 Nacimiento:   %s\n", autor.FechaNac)
	fmt.Printf("📍 Lugar:        %s\n", autor.LugarNac)
	fmt.Printf("🌐 URL perfil:   %s\n", autor.URLPerfil)
	fmt.Printf("📖 Descripción:  %s\n", desc)
	fmt.Println("\n✅ Ejercicio 3 completado!")
}