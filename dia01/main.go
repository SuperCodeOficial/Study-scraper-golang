package main

import (
	"fmt"       // Para imprimir texto
	"io/ioutil" // Para leer el cuerpo de la respuesta
	"net/http"  // Para manejar las peticiones HTTP
	"os"        // Para manejar errores del sistema
)

func main() {
    // 1. Definimos la URL de nuestra "víctima"
    url := "https://webscraper.io/test-sites/e-commerce/allinone"

    fmt.Println("🤖 Iniciando petición a:", url)

    // 2. Hacemos la petición GET
    resp, err := http.Get(url)

    // 3. Manejo de errores básico
    if err != nil {
        fmt.Println("Error al conectar:", err)
        os.Exit(1)
    }
    defer resp.Body.Close() // Cerramos la conexión al terminar

    // 4. Verificamos si la página existe
    if resp.StatusCode != 200 {
        fmt.Println("La web respondió con error:", resp.StatusCode)
        return
    }

    // 5. Leemos el cuerpo de la respuesta
    bodyBytes, _ := ioutil.ReadAll(resp.Body)

    // 6. Convertimos los bytes a String
    htmlContent := string(bodyBytes)

    // 7. Imprimimos los primeros 500 caracteres
    fmt.Println("\n✅ ¡Éxito! Aquí tienes un fragmento del HTML:\n")
    fmt.Println(htmlContent[:500])
    fmt.Println("\n... (y mucho más código HTML)")
}
