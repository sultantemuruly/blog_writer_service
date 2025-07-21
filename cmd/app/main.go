package main

import (
    "os"
    "log"
    "net/http"
    "time"
    "context"

    "github.com/joho/godotenv"
    // "github.com/jackc/pgx/v5"

    "github.com/sultantemuruly/blog_writer_service/internal/db"
	"github.com/sultantemuruly/blog_writer_service/internal/routes"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, continuing")
    }

    ctx := context.Background()

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("NEON_DATABASE_URL environment variable is not set")
        log.Println("Please set the NEON_DATABASE_URL environment variable to your database connection string.")
        return
    }

    conn, err := db.Connect(ctx, dsn)
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
        log.Println("Ensure your database is running and the connection string is correct.")
    }
    defer conn.Close(ctx)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, conn)

    client := &http.Client{Timeout: 40 * time.Second}
    go func () {
        ticker := time.NewTicker(60 * time.Second)
        defer ticker.Stop()

        if err := WriteBlog(client); err != nil {
            log.Println("Error writing blog:", err)
            log.Println("Retrying in 60 seconds...")
            ticker.Reset(60 * time.Second)
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
