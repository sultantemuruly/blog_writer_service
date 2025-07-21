package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/llm_response", llmHandler)
	mux.HandleFunc("/blogs", postBlog(conn))
}
