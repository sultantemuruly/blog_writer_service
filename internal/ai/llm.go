package ai

import (
    "fmt"
    "os"

    "github.com/tmc/langchaingo/llms"
    "github.com/tmc/langchaingo/llms/openai"
)

func NewLLM() (llms.LLM, error) {
    key := os.Getenv("AZURE_OPENAI_API_KEY")
    if key == "" {
        return nil, fmt.Errorf("AZURE_OPENAI_API_KEY not set")
    }
    base := os.Getenv("AZURE_OPENAI_API_BASE_URL")
    if base == "" {
        return nil, fmt.Errorf("AZURE_OPENAI_API_BASE_URL not set")
    }
    ver := os.Getenv("AZURE_OPENAI_API_VERSION")
    if ver == "" {
        return nil, fmt.Errorf("AZURE_OPENAI_API_VERSION not set")
    }
    dep := os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME")
    if dep == "" {
        return nil, fmt.Errorf("AZURE_OPENAI_DEPLOYMENT_NAME not set")
    }

    return openai.New(
        openai.WithToken(key),
        openai.WithBaseURL(base),
        openai.WithAPIVersion(ver),
        openai.WithAPIType(openai.APITypeAzure),
        openai.WithModel(dep),
    )
}
