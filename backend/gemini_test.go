package backend

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeminiClient_SummarizeText_Sanitization(t *testing.T) {
	var receivedPath string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"candidates": [
				{
					"content": {
						"parts": [
							{
								"text": "Mock response content"
							}
						],
						"role": "model"
					},
					"finishReason": "STOP"
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewGeminiClient("test-key")
	client.BaseURL = server.URL

	testCases := []struct {
		name         string
		modelInput   string
		expectedPath string
	}{
		{
			name:         "WithoutPrefix",
			modelInput:   "gemini-3.1-flash-lite",
			expectedPath: "/v1beta/models/gemini-3.1-flash-lite:generateContent",
		},
		{
			name:         "WithPrefix",
			modelInput:   "models/gemini-3.1-flash-lite",
			expectedPath: "/v1beta/models/gemini-3.1-flash-lite:generateContent",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			receivedPath = ""
			summary, _, err := client.SummarizeText(tc.modelInput, "System prompt", "Hello")
			if err != nil {
				t.Fatalf("SummarizeText failed: %v", err)
			}
			if summary != "Mock response content" {
				t.Errorf("Expected summary 'Mock response content', got '%s'", summary)
			}
			if receivedPath != tc.expectedPath {
				t.Errorf("Expected request path to be '%s', got '%s'", tc.expectedPath, receivedPath)
			}
		})
	}
}
