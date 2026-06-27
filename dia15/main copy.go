package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

// ============================================================
// CONCEPTO 1: Struct — agrupa los datos de una cita
// ============================================================
type Cita struct {
	Texto  string
	Autor  string
	Tag    string
}

func main() {

	// ============================================================
	// CONCEPTO 2: ExecAllocator — configura CÓMO lanza Chrome
	// Sin esto: Chrome corre invisible (headless)
	// Con headless:false — Chrome se abre visualmente
	// ============================================================
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),          // 👁️ Ver Chrome abrirse
		chromedp.Flag("disable-gpu", true),        // Estabilidad en Windows
		chromedp.Flag("no-sandbox", true),         // Necesario en algunos sistemas
		chromedp.WindowSize(1280, 800),            // Tamaño de ventana
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// ============================================================
	// CONCEPTO 3: NewContext — crea la "sesión" del navegador
	// Es el punto de control de TODO lo que hace Chrome
	// ============================================================
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf), // Ver logs internos de chromedp
	)
	defer cancel()

	// ============================================================
	// CONCEPTO 4: WithTimeout — límite de tiempo global
	// Si en 60s no termina → cancela todo automáticamente
	// ============================================================
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// ============================================================
	// CONCEPTO 5: Variables receptoras
	// chromedp escribirá los resultados aquí usando & (punteros)
	// ============================================================
	var urlActual    string
	var tituloPagina string
	var cita1        Cita
	var cita2        Cita
	var cita3        Cita
	var totalCitas   string
	var htmlCita     string

	// ============================================================
	// CONCEPTO 6: chromedp.Run() — ejecuta acciones en secuencia
	// Cada acción se ejecuta UNA POR UNA dentro del Chrome real
	// ============================================================
	fmt.Println("🌐 Lanzando Chrome...")

	err := chromedp.Run(ctx,

		// --- ACCIÓN 1: Navegar ---
		// Chrome abre la URL como si escribieras en la barra
		chromedp.Navigate("https://quotes.toscrape.com/"),

		// --- ACCIÓN 2: WaitVisible ---
		// ESPERA hasta que el elemento aparezca en pantalla
		// Sin esto, intentaríamos leer antes de que cargue
		chromedp.WaitVisible(`.quote`, chromedp.ByQuery),

		// --- ACCIÓN 3: Pausa humana ---
		// Damos tiempo extra para que todo renderice
		chromedp.Sleep(2*time.Second),

		// --- ACCIÓN 4: Capturar URL actual ---
		chromedp.Location(&urlActual),

		// --- ACCIÓN 5: Título de la página (pestaña del navegador) ---
		chromedp.Title(&tituloPagina),

		// --- ACCIÓN 6: Extraer texto con selector CSS ---
		// .quote:nth-child(1) = primera cita
		// .text = el span con el texto de la cita
		chromedp.Text(`.quote:nth-child(1) .text`, &cita1.Texto, chromedp.ByQuery),
		chromedp.Text(`.quote:nth-child(1) .author`, &cita1.Autor, chromedp.ByQuery),
		chromedp.Text(`.quote:nth-child(1) .tag`, &cita1.Tag, chromedp.ByQuery),

		// --- ACCIÓN 7: Segunda cita ---
		chromedp.Text(`.quote:nth-child(2) .text`, &cita2.Texto, chromedp.ByQuery),
		chromedp.Text(`.quote:nth-child(2) .author`, &cita2.Autor, chromedp.ByQuery),
		chromedp.Text(`.quote:nth-child(2) .tag`, &cita2.Tag, chromedp.ByQuery),

		chromedp.Text(`.quote:nth-child(3) .text`, &cita3.Texto, chromedp.ByQuery),
		chromedp.Text(`.quote:nth-child(3) .author`, &cita3.Autor, chromedp.ByQuery),
		chromedp.Text(`.quote:nth-child(3) .tag`, &cita3.Tag, chromedp.ByQuery),

		// --- ACCIÓN 8: Contar elementos con JavaScript ---
		// Evaluate() ejecuta JS directo en el navegador
		// Aquí contamos cuántas citas hay en la página
		chromedp.Evaluate(
			`document.querySelectorAll('.quote').length.toString()`,
			&totalCitas,
		),

		// --- ACCIÓN 9: Obtener HTML de un elemento ---
		// OuterHTML nos da el HTML completo de un nodo
		chromedp.OuterHTML(`.quote:nth-child(1)`, &htmlCita, chromedp.ByQuery),
	)

	// ============================================================
	// CONCEPTO 7: Manejo de errores
	// Si CUALQUIER acción falla, err tendrá el motivo
	// ============================================================
	if err != nil {
		log.Fatal("❌ Error en chromedp:", err)
	}

	// ============================================================
	// IMPRIMIR RESULTADOS
	// ============================================================
	fmt.Println("\n╔══════════════════════════════════════╗")
	fmt.Println("║      RESULTADO DÍA 15 — chromedp     ║")
	fmt.Println("╚══════════════════════════════════════╝")

	fmt.Printf("\n🌐 URL visitada:    %s\n", urlActual)
	fmt.Printf("📄 Título pestaña:  %s\n", tituloPagina)
	fmt.Printf("📊 Citas en página: %s\n", totalCitas)

	fmt.Println("\n--- CITA #1 ---")
	fmt.Printf("💬 Texto:  %s\n", cita1.Texto)
	fmt.Printf("✍️  Autor:  %s\n", cita1.Autor)
	fmt.Printf("🏷️  Tag:    %s\n", cita1.Tag)

	fmt.Println("\n--- CITA #2 ---")
	fmt.Printf("💬 Texto:  %s\n", cita2.Texto)
	fmt.Printf("✍️  Autor:  %s\n", cita2.Autor)
	fmt.Printf("🏷️  Tag:    %s\n", cita2.Tag)

	fmt.Println("\n--- CITA #3 ---")
	fmt.Printf("💬 Texto:  %s\n", cita3.Texto)
	fmt.Printf("✍️  Autor:  %s\n", cita3.Autor)
	fmt.Printf("🏷️  Tag:    %s\n", cita3.Tag)

	fmt.Println("\n--- HTML REAL DEL ELEMENTO (primeros 300 chars) ---")
	if len(htmlCita) > 300 {
		fmt.Println(htmlCita[:300] + "...")
	} else {
		fmt.Println(htmlCita)
	}

	fmt.Println("\n✅ Día 15 completado con éxito!")
}
