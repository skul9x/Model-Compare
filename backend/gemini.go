package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GeminiContentPart represents a single part of content (usually text)
type GeminiContentPart struct {
	Text string `json:"text"`
}

// GeminiContent represents the content structure for Gemini API
type GeminiContent struct {
	Parts []GeminiContentPart `json:"parts"`
}

// GeminiSystemInstruction represents the system instruction structure for Gemini API
type GeminiSystemInstruction struct {
	Parts []GeminiContentPart `json:"parts"`
}

// GeminiRequest represents the request body for Gemini API generateContent
type GeminiRequest struct {
	Contents          []GeminiContent          `json:"contents"`
	SystemInstruction *GeminiSystemInstruction `json:"systemInstruction,omitempty"`
}

// GeminiResponse represents the response body from Gemini API
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

// GeminiClient manages the API Key and coordinates requests to Gemini API
type GeminiClient struct {
	ApiKey  string
	BaseURL string
}

// NewGeminiClient creates a new GeminiClient
func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		ApiKey:  apiKey,
		BaseURL: "https://generativelanguage.googleapis.com",
	}
}

// SummarizeText calls the Gemini API to summarize text with optional system prompt
// It returns the summarized text, response time in milliseconds, and any error.
func (c *GeminiClient) SummarizeText(modelName string, systemPrompt string, text string) (string, int64, error) {
	startTime := time.Now()

	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com"
	}
	cleanedModelName := strings.TrimPrefix(modelName, "models/")
	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", baseURL, cleanedModelName, c.ApiKey)

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiContentPart{
					{Text: text},
				},
			},
		},
	}

	if systemPrompt != "" {
		reqBody.SystemInstruction = &GeminiSystemInstruction{
			Parts: []GeminiContentPart{
				{Text: systemPrompt},
			},
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		duration := time.Since(startTime).Milliseconds()
		return "", duration, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		duration := time.Since(startTime).Milliseconds()
		return "", duration, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(startTime).Milliseconds()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", duration, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp GeminiResponse
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Error != nil {
			return "", duration, fmt.Errorf("gemini API error (status %d): [%s] %s", resp.StatusCode, errResp.Error.Status, errResp.Error.Message)
		}
		return "", duration, fmt.Errorf("gemini API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return "", duration, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if geminiResp.Error != nil {
		return "", duration, fmt.Errorf("gemini API error: [%s] %s", geminiResp.Error.Status, geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", duration, fmt.Errorf("received empty candidates list from Gemini API")
	}

	outputText := geminiResp.Candidates[0].Content.Parts[0].Text
	return outputText, duration, nil
}
