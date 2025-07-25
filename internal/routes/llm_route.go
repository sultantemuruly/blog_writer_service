package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sultantemuruly/blog_writer_service/internal/ai"
	"github.com/tmc/langchaingo/llms"
)

// type llmResponse struct {
// 	Response string `json:"response"`
// }

type BlogRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type blogResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func llmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prompt := "Write a blog post about something that is related to emails. Like the history of email, how to properly write an email, etc. you should respond using title, and html body which will be the blog itself"

	// Prepare the system prompt for the LLM
	systemPrompt := fmt.Sprintf(
		`You are a blog-writing assistant. 
        Please respond with *only* a JSON object with exactly two string fields: 
        "title" – a one-line blog post title, 
        "content" – the full body text. 
        Do not wrap the JSON in markdown or text. 
        User request: %s`,
		prompt,
	)

	ctx := r.Context()
	llm, err := ai.NewLLM()
	if err != nil {
		log.Println("Error creating LLM:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Generate the response from the LLM
	raw, err := llms.GenerateFromSinglePrompt(ctx, llm, systemPrompt)
	if err != nil {
		log.Println("Error generating completion:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Verify the response format
	var resp blogResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		log.Println("LLM returned invalid JSON:", raw, err)
		http.Error(w, "failed to parse LLM output", http.StatusInternalServerError)
		return
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Error encoding JSON response:", err)
	}
}
