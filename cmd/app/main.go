package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	// "github.com/jackc/pgx/v5"

	"github.com/sultantemuruly/blog_writer_service/internal/db"
	"github.com/sultantemuruly/blog_writer_service/internal/routes"
)

func main() {
	if os.Getenv("DOCKER_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, continuing")
		}
	}

	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = os.Getenv("NEON_DATABASE_URL")
	}
	if dsn == "" {
		log.Fatal("DATABASE_URL or NEON_DATABASE_URL environment variable is not set")
		log.Println("Please set the DATABASE_URL or NEON_DATABASE_URL environment variable to your database connection string.")
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
	go func() {
		// Wait for server to be ready before starting scheduled tasks
		time.Sleep(5 * time.Second)

		for {
			// Try up to 10 times, with 1 minute between attempts
			for i := 0; i < 10; i++ {
				if err := WriteBlog(client); err != nil {
					log.Println("Error writing blog:", err)
					if i < 9 {
						log.Println("Retrying in 1 minute...")
						time.Sleep(1 * time.Minute)
					}
				} else {
					break // Success, exit retry loop
				}
			}
			// Wait 24 hours before next scheduled run
			time.Sleep(24 * time.Hour)
		}
	}()

	addr := ":8080"
	log.Printf("Server listening on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func WriteBlog(client *http.Client) error {
	resp, err := client.Get("http://localhost:8080/llm_response")
	if err != nil {
		log.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	log.Println("🕒 Scheduled GET /llm_response", resp.Status)

	// Parse the response into blogResponse struct
	var blog struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&blog); err != nil {
		log.Println("Error decoding /llm_response JSON:", err)
		return err
	}

	// Marshal the blog struct to JSON for POST
	blogJSON, err := json.Marshal(blog)
	if err != nil {
		log.Println("Error marshaling blog to JSON:", err)
		return err
	}

	// POST to /blogs
	postResp, err := client.Post("http://localhost:8080/blogs", "application/json", bytes.NewReader(blogJSON))
	if err != nil {
		log.Println("Error posting to /blogs:", err)
		return err
	}
	defer postResp.Body.Close()

	log.Println("Blog POST /blogs", postResp.Status)
	return nil
}
