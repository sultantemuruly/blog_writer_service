package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/sultantemuruly/blog_writer_service/internal/ai"
	"github.com/tmc/langchaingo/llms"
)

type BlogRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type blogResponse struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Blog topic categories with diverse subtopics
var blogTopics = map[string][]string{
	"Technology": {
		"Artificial Intelligence and Machine Learning",
		"Cybersecurity and Privacy",
		"Cloud Computing and DevOps",
		"Mobile App Development",
		"Web Development Trends",
		"Data Science and Analytics",
		"Internet of Things (IoT)",
		"Blockchain and Cryptocurrency",
		"Virtual Reality and Augmented Reality",
		"Software Architecture Patterns",
	},
	"Productivity": {
		"Time Management Techniques",
		"Remote Work Best Practices",
		"Project Management Strategies",
		"Goal Setting and Achievement",
		"Work-Life Balance",
		"Stress Management",
		"Digital Organization",
		"Communication Skills",
		"Leadership Development",
		"Personal Development",
	},
	"Business": {
		"Startup Strategies",
		"Digital Marketing",
		"Customer Experience",
		"Sales Techniques",
		"Financial Management",
		"Team Building",
		"Innovation and Creativity",
		"Market Research",
		"Brand Building",
		"Entrepreneurship",
	},
	"Communication": {
		"Email Writing Best Practices",
		"Professional Communication",
		"Public Speaking",
		"Negotiation Skills",
		"Cross-cultural Communication",
		"Conflict Resolution",
		"Presentation Skills",
		"Networking Strategies",
		"Social Media Communication",
		"Technical Writing",
	},
	"Health & Wellness": {
		"Mental Health Awareness",
		"Physical Fitness",
		"Nutrition and Diet",
		"Sleep Optimization",
		"Mindfulness and Meditation",
		"Workplace Wellness",
		"Digital Detox",
		"Stress Relief Techniques",
		"Healthy Habits",
		"Work-Life Integration",
	},
}

// Writing styles for variety
var writingStyles = []string{
	"Professional and authoritative",
	"Conversational and engaging",
	"Educational and tutorial-style",
	"Inspirational and motivational",
	"Analytical and data-driven",
	"Storytelling and narrative",
	"Problem-solving and solution-focused",
	"Trend analysis and future-focused",
}

func getRandomTopic() (string, string) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Select random category
	categories := make([]string, 0, len(blogTopics))
	for category := range blogTopics {
		categories = append(categories, category)
	}

	selectedCategory := categories[rand.Intn(len(categories))]
	subtopics := blogTopics[selectedCategory]
	selectedSubtopic := subtopics[rand.Intn(len(subtopics))]

	return selectedCategory, selectedSubtopic
}

func getRandomStyle() string {
	rand.Seed(time.Now().UnixNano())
	return writingStyles[rand.Intn(len(writingStyles))]
}

func llmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get random topic and style
	category, subtopic := getRandomTopic()
	style := getRandomStyle()

	// Create a dynamic, detailed prompt
	prompt := fmt.Sprintf(
		`Write a comprehensive, high-quality blog post about "%s" within the "%s" category. 
		
		Requirements:
		- Title: Create an engaging, SEO-friendly title (50-60 characters)
		- Content: Write a well-structured blog post (800-1200 words) with:
			* Compelling introduction that hooks the reader
			* Clear sections with subheadings
			* Practical tips, examples, or actionable advice
			* Engaging conclusion with call-to-action
			* Use bullet points and numbered lists where appropriate
			* Include relevant statistics or facts when possible
		
		Writing Style: %s
		
		Focus on providing genuine value to readers. Make the content informative, engaging, and actionable.`,
		subtopic, category, style,
	)

	// Enhanced system prompt for better AI responses
	systemPrompt := fmt.Sprintf(
		`You are an expert content writer and blogger with deep knowledge across multiple domains. 
		
		Your task is to create high-quality, engaging blog content that provides real value to readers.
		
		IMPORTANT INSTRUCTIONS:
		1. Respond with ONLY a valid JSON object containing exactly two fields:
		   - "title": A compelling, SEO-optimized title (50-60 characters)
		   - "content": A well-structured HTML blog post (800-1200 words)
		
		2. Content Requirements:
		   - Use proper HTML formatting (h2, h3, p, ul, ol, strong, em)
		   - Include engaging subheadings
		   - Add bullet points and numbered lists for readability
		   - Use bold and italic text for emphasis
		   - Include practical examples and actionable tips
		   - End with a compelling conclusion
		
		3. Quality Standards:
		   - Write in a %s tone
		   - Ensure content is informative and valuable
		   - Make it engaging and easy to read
		   - Include relevant insights and practical advice
		
		4. Format the content as clean HTML with proper tags
		
		User Request: %s
		
		Remember: Return ONLY the JSON object, no additional text or markdown formatting.`,
		style, prompt,
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

	// Log the generated topic for debugging
	log.Printf("Generated blog about: %s - %s (Style: %s)", category, subtopic, style)

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
