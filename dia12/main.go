package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
	_ "modernc.org/sqlite" // <--- CAMBIO 1: Usamos la librería "modernc"
)

type Product struct {
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Stock    int64   `json:"stock"`
}

type APIResponse struct {
	Products []Product `json:"products"`
}

func main() {
	// --- A. CONFIGURACIÓN DE LA BASE DE DATOS ---
	
	// CAMBIO 2: El nombre del driver ahora es "sqlite", no "sqlite3"
	db, err := sql.Open("sqlite", "./inventario.db") 
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Creamos la tabla si no existe
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS productos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		price REAL,
		category TEXT,
		stock INTEGER
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("Error creando tabla: %q: %s\n", err, sqlStmt)
		return
	}
	fmt.Println("💾 Base de datos lista: inventario.db")

	// --- B. CONFIGURACIÓN DEL SCRAPER ---
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		var respuesta APIResponse
		if err := json.Unmarshal(r.Body, &respuesta); err != nil {
			log.Println("Error JSON:", err)
			return
		}

		fmt.Printf("📥 Procesando %d productos...\n", len(respuesta.Products))

		// --- C. GUARDADO EN BASE DE DATOS ---
		stmt, err := db.Prepare("INSERT INTO productos(title, price, category, stock) VALUES(?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		for _, p := range respuesta.Products {
			_, err = stmt.Exec(p.Title, p.Price, p.Category, p.Stock)
			if err != nil {
				log.Println("Error al guardar:", err)
			} else {
				fmt.Printf("   ✅ Guardado: %s\n", p.Title)
			}
		}
	})

	fmt.Println("--- Iniciando Scraper con BD (Pure Go) ---")
	c.Visit("https://dummyjson.com/products?limit=10")
}
