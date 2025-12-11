package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
)

// OpenAILLM implements the llms.LLM interface using OpenAI as the backend
type OpenAILLM struct {
	client           *openai.Client
	model            string
	temperature      *float64
	maxTokens        *int
	topP             *float64
	presencePenalty  *float64
	frequencyPenalty *float64
}

// NewOpenAILLM creates a new instance of OpenAILLM using OpenAI as the backend
// The API key is read from the OPENAI_API_KEY environment variable
func NewOpenAILLM(opts ...Option) (*OpenAILLM, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	llm := &OpenAILLM{
		client: openai.NewClient(apiKey),
		model:  "gpt-3.5-turbo-instruct", // Default model
	}

	for _, opt := range opts {
		opt(llm)
	}

	return llm, nil
}

// Call implements the llms.LLM interface
func (o *OpenAILLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	// Apply any call options
	opts := &llms.CallOptions{}
	for _, opt := range options {
		opt(opts)
	}
	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model: o.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   1000, // Default value, can be overridden by options
		Temperature: 0.7,  // Default value, can be overridden by options
	}

	// Apply options from the OpenAILLM instance
	if o.temperature != nil {
		req.Temperature = float32(*o.temperature)
	}
	// Apply temperature from call options if provided
	if opts.Temperature != 0 {
		req.Temperature = float32(opts.Temperature)
	}
	if o.maxTokens != nil {
		req.MaxTokens = *o.maxTokens
	}
	if o.topP != nil {
		req.TopP = float32(*o.topP)
	}
	if o.presencePenalty != nil {
		req.PresencePenalty = float32(*o.presencePenalty)
	}
	if o.frequencyPenalty != nil {
		req.FrequencyPenalty = float32(*o.frequencyPenalty)
	}

	// Make the API call using ChatCompletion
	resp, err := o.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	// Return the content of the first choice's message
	return resp.Choices[0].Message.Content, nil
}

// GenerateResult represents a single generation result
type GenerateResult struct {
	Text string
}

// Generate processes multiple prompts and returns results
func (o *OpenAILLM) Generate(ctx context.Context, prompts []string, options ...llms.CallOption) ([]*GenerateResult, error) {
	results := make([]*GenerateResult, 0, len(prompts))

	for _, prompt := range prompts {
		text, err := o.Call(ctx, prompt, options...)
		if err != nil {
			return nil, err
		}

		results = append(results, &GenerateResult{
			Text: text,
		})
	}

	return results, nil
}

// GetNumTokens implements the llms.LLM interface
func (o *OpenAILLM) GetNumTokens(text string) int {
	// This is a rough approximation
	// For more accurate token counting, you might want to use the tiktoken-go package
	return len(text) / 4
}

// GetType returns the type of the LLM
func (o *OpenAILLM) GetType() string {
	return "OpenAI"
}

// Option is a function that configures an OpenAILLM
// This pattern allows for flexible configuration of the LLM
// Example: llm.NewOpenAILLM(llm.WithModel("gpt-4"), llm.WithTemperature(0.7))
type Option func(*OpenAILLM)

// WithModel sets the model to use (e.g., "gpt-3.5-turbo-instruct", "gpt-4")
func WithModel(model string) Option {
	return func(o *OpenAILLM) {
		o.model = model
	}
}

// WithTemperature sets the temperature parameter (0.0 to 2.0)
func WithTemperature(temp float64) Option {
	return func(o *OpenAILLM) {
		o.temperature = &temp
	}
}

// WithMaxTokens sets the maximum number of tokens to generate
func WithMaxTokens(maxTokens int) Option {
	return func(o *OpenAILLM) {
		o.maxTokens = &maxTokens
	}
}

// WithTopP sets the nucleus sampling parameter (0.0 to 1.0)
func WithTopP(topP float64) Option {
	return func(o *OpenAILLM) {
		o.topP = &topP
	}
}

// WithPresencePenalty sets the presence penalty (-2.0 to 2.0)
func WithPresencePenalty(penalty float64) Option {
	return func(o *OpenAILLM) {
		o.presencePenalty = &penalty
	}
}

// WithFrequencyPenalty sets the frequency penalty (-2.0 to 2.0)
func WithFrequencyPenalty(penalty float64) Option {
	return func(o *OpenAILLM) {
		o.frequencyPenalty = &penalty
	}
}
