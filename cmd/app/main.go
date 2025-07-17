package main

import (
    "log"
    "net/http"
    "time"

    "github.com/joho/godotenv"

	"github.com/sultantemuruly/blog_writer_service/internal/routes"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, continuing")
    }

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

    client := &http.Client{Timeout: 40 * time.Second}
    go func () {
        ticker := time.NewTicker(60 * time.Second)
        defer ticker.Stop()

        if err := WriteBlog(client); err != nil {
            log.Println("Error writing blog:", err)
        }

        for range ticker.C {
            if err := WriteBlog(client); err != nil {
                log.Println("Error writing blog:", err)
            }
        }
    }()

	addr := ":8080"
    log.Printf("Server listening on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func WriteBlog(client *http.Client) error {
    resp, err := client.Get("http://localhost:8080/llm_response?prompt=Write%20a%20blog%20post%20about%20Go%20programming")
    if err != nil {
        log.Println("Error making request:", err)
        return err
    }
    defer resp.Body.Close()

    log.Println("🕒 Scheduled GET /llm_response", resp.Status)
    return nil
}
