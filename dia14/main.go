package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
	_ "modernc.org/sqlite" // <--- 1. CAMBIO AQUÍ: Usamos el driver puro de Go
)

type Noticia struct {
	Fuente string
	Titulo string
	Link   string
	Fecha  string
}

func main() {
	fmt.Println("📰 Iniciando Monitor de Noticias Multihilo...")

	// 2. CAMBIO AQUÍ: El nombre del driver ahora es "sqlite" (sin el 3)
	db, err := sql.Open("sqlite", "./noticias.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlTabla := `
	CREATE TABLE IF NOT EXISTS noticias (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fuente TEXT,
		titulo TEXT,
		link TEXT,
		fecha DATETIME
	);`
	_, err = db.Exec(sqlTabla)
	if err != nil {
		log.Fatal("Error creando tabla:", err)
	}

	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*books.toscrape.com*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*scrapeme.live*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	// SITIO A
	c.OnHTML("article.product_pod h3 a", func(e *colly.HTMLElement) {
		titulo := e.Attr("title")
		link := e.Request.AbsoluteURL(e.Attr("href"))
		guardarEnBD(db, "Librería News", titulo, link)
	})

	// SITIO B
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		titulo := e.ChildText("h2")
		link := e.ChildAttr("a", "href")
		guardarEnBD(db, "PokeMundo", titulo, link)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("🚀 Visitando:", r.URL)
	})

	sitios := []string{
		"http://books.toscrape.com/catalogue/category/books/travel_2/index.html",
		"https://scrapeme.live/shop/",
	}

	inicio := time.Now()

	for _, url := range sitios {
		c.Visit(url)
	}

	c.Wait()

	fmt.Printf("\n✅ Proceso terminado en %s.\n", time.Since(inicio))
	fmt.Println("📂 Revisa tu archivo 'noticias.db'")
}

func guardarEnBD(db *sql.DB, fuente, titulo, link string) {
	stmt, err := db.Prepare("INSERT INTO noticias(fuente, titulo, link, fecha) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparando SQL:", err)
		return
	}
	defer stmt.Close()

	fechaActual := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(fuente, titulo, link, fechaActual)
	if err != nil {
		log.Println("Error insertando:", err)
	} else {
		fmt.Print(".")
	}
}
