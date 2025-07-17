package main

import (
    "log"
    "net/http"

    "github.com/joho/godotenv"

	"github.com/sultantemuruly/blog_writer_service/internal/routes"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, continuing")
    }

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	addr := ":8080"
    log.Printf("Server listening on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
