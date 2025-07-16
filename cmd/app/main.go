package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/sultantemuruly/blog_writer_service/internal/ai"
)

func main() {
	
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, relying on real ENV vars")
    }

	ctx := context.Background()

	llm, err := ai.NewLLM()
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error creating LLM:", err)
		return
	}

	completion, err := llms.GenerateFromSinglePrompt(ctx,
		llm,
		"Tell me a joke please!",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion)
}