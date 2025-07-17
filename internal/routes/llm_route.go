package routes

import (
	"context"
	"net/http"
	"encoding/json"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/sultantemuruly/blog_writer_service/internal/ai"
)

type llmResponse struct {
    Response string `json:"response"`
}

func llmHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    prompt := r.URL.Query().Get("prompt")
    if prompt == "" {
        http.Error(w, "`prompt` query parameter is required", http.StatusBadRequest)
        return
    }

    ctx := context.Background()
    llm, err := ai.NewLLM()
    if err != nil {
        log.Println("Error creating LLM:", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
    if err != nil {
        log.Println("Error generating completion:", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    json.NewEncoder(w).Encode(llmResponse{Response: completion})
}
