package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Definimos la carpeta como una constante para usarla en todo el programa
const carpetaSalida = "reseñas_tablets"

func main() {
	// 1. Creamos la carpeta al inicio
	os.MkdirAll(carpetaSalida, os.ModePerm)

	c := colly.NewCollector()

	// ---------------------------------------------------------
	// PASO 1 y 2: NAVEGACIÓN (Sidebar)
	// ---------------------------------------------------------
	// Usamos un solo bloque para manejar el menú lateral.
	// Buscamos enlaces dentro de .sidebar-nav
	c.OnHTML(".sidebar-nav a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		texto := strings.TrimSpace(e.Text)

		// Si es Computers, entramos
		if texto == "Computers" {
			fmt.Println("1. 📂 Entrando a categoría:", texto)
			e.Request.Visit(link)
		}

		// Si es Tablets, entramos
		if texto == "Tablets" {
			fmt.Println("2. 📂 Entrando a sub-categoría:", texto)
			e.Request.Visit(link)
		}
	})

	// ---------------------------------------------------------
	// PASO 3: DE LA LISTA AL DETALLE
	// ---------------------------------------------------------
	c.OnHTML("div.thumbnail", func(e *colly.HTMLElement) {
		// Verificamos la URL para asegurar que estamos en la sección de Tablets
		urlActual := e.Request.URL.String()

		if strings.Contains(urlActual, "/tablets") {
			// Buscamos el enlace del título dentro del thumbnail
			linkProducto := e.ChildAttr("a.title", "href")
			
			fmt.Println("   🔍 Tablet detectada, entrando al detalle:", linkProducto)
			
			// Visitamos el producto (esto disparará el PASO 4)
			e.Request.Visit(linkProducto)
		}
	})

	// ---------------------------------------------------------
	// PASO 4: EXTRACCIÓN Y GUARDADO (Detalle del producto)
	// ---------------------------------------------------------
	c.OnHTML("div.caption", func(e *colly.HTMLElement) {
		// IMPORTANTE: Este selector existe tanto en la lista como en el detalle.
		// Debemos asegurarnos de NO estar en la lista para no guardar basura.
		urlActual := e.Request.URL.String()
		
		// Si la URL termina en "tablets", significa que estamos en la lista, NO guardar.
		if strings.HasSuffix(urlActual, "/tablets") {
			return 
		}

		// Si pasamos el filtro, extraemos los datos
		nombre := e.ChildText("h4:nth-of-type(2)")
		descripcion := e.ChildText("p.description")

		// Verificación final de que hay datos
		if nombre != "" && descripcion != "" {
			guardarArchivo(nombre, descripcion)
		}
	})

	// Inicio del robot
	fmt.Println("🤖 Iniciando robot explorador...")
	c.Visit("https://webscraper.io/test-sites/e-commerce/allinone")
}

// Función auxiliar para guardar el archivo
func guardarArchivo(nombre, contenido string) {
	// Limpieza del nombre para que sea un archivo válido
	nombreLimpio := strings.ReplaceAll(nombre, " ", "_")
	nombreLimpio = strings.ReplaceAll(nombreLimpio, "/", "-")
	
	// Construimos la ruta completa: "reseñas_tablets/Lenovo_IdeaTab.txt"
	rutaCompleta := carpetaSalida + "/" + nombreLimpio + ".txt"

	f, err := os.Create(rutaCompleta)
	if err != nil {
		fmt.Println("Error creando archivo:", err)
		return
	}
	defer f.Close()

	f.WriteString(contenido)
	fmt.Printf("   💾 Guardado exitoso: %s\n", nombreLimpio)
}
