package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var GeminiBaseURL = "https://generativelanguage.googleapis.com"

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content ResponseContent `json:"content"`
}

type ResponseContent struct {
	Parts []Part `json:"parts"`
}

// CallGemini calls the Google Gemini API to generate content from a prompt
func CallGemini(ctx context.Context, apiKey, model, prompt string) (string, error) {
	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", GeminiBaseURL, model, apiKey)

	reqPayload := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Gemini request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request for Gemini: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Production grade: set request timeouts
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http call to Gemini failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API returned bad status code: %d", resp.StatusCode)
	}

	var respPayload GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&respPayload); err != nil {
		return "", fmt.Errorf("failed to decode Gemini response: %w", err)
	}

	if len(respPayload.Candidates) == 0 ||
		len(respPayload.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("received empty response candidates from Gemini")
	}

	return respPayload.Candidates[0].Content.Parts[0].Text, nil
}
