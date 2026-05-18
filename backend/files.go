package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AppConfig stores the user's settings and preferences
type AppConfig struct {
	ApiKey       string   `json:"apiKey"`
	Models       []string `json:"models"`
	SystemPrompt string   `json:"systemPrompt"`
	OutputDir    string   `json:"outputDir"`
}

// configPathOverride is used to redirect configuration storage during unit tests
var configPathOverride string

// getConfigFilepath returns the path to the configuration file, respecting overrides
func getConfigFilepath() (string, error) {
	if configPathOverride != "" {
		return configPathOverride, nil
	}
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	appConfigDir := filepath.Join(userConfigDir, "modeltester")
	return filepath.Join(appConfigDir, "config.json"), nil
}

// SaveConfig saves the AppConfig to the platform-specific standard configuration directory
func SaveConfig(cfg AppConfig) error {
	filePath, err := getConfigFilepath()
	if err != nil {
		return err
	}

	appConfigDir := filepath.Dir(filePath)
	if err := os.MkdirAll(appConfigDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadConfig loads the AppConfig from the platform-specific standard configuration directory
func LoadConfig() (AppConfig, error) {
	filePath, err := getConfigFilepath()
	if err != nil {
		return AppConfig{}, err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Return an empty config if the file does not exist yet
		return AppConfig{}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return AppConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return AppConfig{}, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	return cfg, nil
}

// SaveResult writes the summary content to a formatted txt file inside the output directory
func SaveResult(outputDir string, modelIdx int, vbIdx int, content string, responseTimeMs int64) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fileName := fmt.Sprintf("model_%d-vb_%d.txt", modelIdx, vbIdx)
	filePath := filepath.Join(outputDir, fileName)

	// Format response time as seconds with 2 decimal places (e.g. 1.42s)
	responseTimeSeconds := float64(responseTimeMs) / 1000.0
	fileContent := fmt.Sprintf("Response Time: %.2fs (%dms)\n---\n%s", responseTimeSeconds, responseTimeMs, content)

	if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
		return fmt.Errorf("failed to write result file: %w", err)
	}

	return nil
}

// SaveMapping creates the mapping.json file containing mapping between index-based names and real names/texts
func SaveMapping(outputDir string, models []string, texts []string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create mapping directory: %w", err)
	}

	modelsMap := make(map[string]string)
	for i, m := range models {
		modelsMap[fmt.Sprintf("model_%d", i+1)] = m
	}

	textsMap := make(map[string]string)
	for i, t := range texts {
		textsMap[fmt.Sprintf("vb_%d", i+1)] = t
	}

	mapping := map[string]interface{}{
		"models": modelsMap,
		"texts":  textsMap,
	}

	jsonData, err := json.MarshalIndent(mapping, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal mapping data: %w", err)
	}

	mappingFile := filepath.Join(outputDir, "mapping.json")
	if err := os.WriteFile(mappingFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write mapping file: %w", err)
	}

	return nil
}

// SaveModelsTimeFile creates a txt file named model_{time}.txt in the outputDir listing the models mapped to indices
func SaveModelsTimeFile(outputDir string, models []string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	currentTime := time.Now().Format("20060102-150405")
	fileName := fmt.Sprintf("model_%s.txt", currentTime)
	filePath := filepath.Join(outputDir, fileName)

	var sb strings.Builder
	for i, m := range models {
		sb.WriteString(fmt.Sprintf("model_%d: %s\n", i+1, m))
	}

	if err := os.WriteFile(filePath, []byte(sb.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to write models time file: %w", err)
	}

	return filePath, nil
}
