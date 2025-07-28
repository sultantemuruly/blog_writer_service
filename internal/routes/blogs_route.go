package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type blogRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func postBlog(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("Received non-POST request for /blogs")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req blogRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("Error decoding request body:", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if _, err := db.Exec(
			r.Context(),
			"INSERT INTO blogs (title, content) VALUES ($1, $2)",
			req.Title, req.Content,
		); err != nil {
			log.Println("Error inserting blog:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		log.Println("Blog inserted successfully without status")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Blog created successfully without status"))
	}
}
