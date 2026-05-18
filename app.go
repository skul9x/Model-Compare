package main

import (
	"context"
	"fmt"
	"strings"
	"time"
	"modeltester/backend"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	isTest     bool // Set to true in unit/integration tests to bypass Wails runtime calls
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SelectDirectory opens a native directory selection dialog and returns the selected path
func (a *App) SelectDirectory() (string, error) {
	if a.isTest {
		return "/mock/directory", nil
	}
	selectedDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Chọn thư mục lưu kết quả",
	})
	if err != nil {
		return "", fmt.Errorf("failed to open directory dialog: %w", err)
	}
	return selectedDir, nil
}

// LoadConfig loads the application configuration
func (a *App) LoadConfig() (backend.AppConfig, error) {
	return backend.LoadConfig()
}

// SaveConfig saves the application configuration
func (a *App) SaveConfig(cfg backend.AppConfig) error {
	return backend.SaveConfig(cfg)
}

// SummarizeText calls the Gemini API directly to summarize text
func (a *App) SummarizeText(modelName string, systemPrompt string, text string) (string, int64, error) {
	cfg, err := backend.LoadConfig()
	if err != nil {
		return "", 0, fmt.Errorf("failed to load API configuration: %w", err)
	}
	if cfg.ApiKey == "" {
		return "", 0, fmt.Errorf("API key is missing, please save your API key in settings first")
	}

	client := backend.NewGeminiClient(cfg.ApiKey)
	return client.SummarizeText(modelName, systemPrompt, text)
}

// SaveResult saves the summary text output into a model_X-vb_Y.txt file
func (a *App) SaveResult(outputDir string, modelIdx int, vbIdx int, content string, responseTimeMs int64) error {
	return backend.SaveResult(outputDir, modelIdx, vbIdx, content, responseTimeMs)
}

// SaveMapping saves the mapping.json configuration
func (a *App) SaveMapping(outputDir string, models []string, texts []string) error {
	return backend.SaveMapping(outputDir, models, texts)
}

// truncateText helper for Dry Run description
func truncateText(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

// emitEvent safely emits a Wails event, skipping if in testing mode
func (a *App) emitEvent(eventName string, optionalData ...interface{}) {
	if a.isTest || a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, eventName, optionalData...)
}

// StartBatchProcess launches the batch summary execution inside a separate goroutine
func (a *App) StartBatchProcess(cfg backend.AppConfig, texts []string) {
	if a.cancelFunc != nil {
		a.cancelFunc()
	}

	var processCtx context.Context
	processCtx, a.cancelFunc = context.WithCancel(a.ctx)

	go func() {
		defer func() {
			a.cancelFunc = nil
			if cfg.OutputDir != "" && len(cfg.Models) > 0 {
				_, _ = backend.SaveModelsTimeFile(cfg.OutputDir, cfg.Models)
			}
			a.emitEvent("task:completed")
		}()

		for vIdx, text := range texts {
			for mIdx, modelName := range cfg.Models {
				// Check if process has been cancelled
				select {
				case <-processCtx.Done():
					a.emitEvent("task:cancelled")
					return
				default:
				}

				// Emit starting event
				a.emitEvent("task:started", map[string]interface{}{
					"modelIdx": mIdx,
					"vbIdx":    vIdx,
				})

				var resultText string
				var responseTimeMs int64
				var err error

				trimmedText := strings.TrimSpace(text)
				if trimmedText == "" {
					err = fmt.Errorf("văn bản trống hoặc chỉ chứa khoảng trắng")
				} else if cfg.ApiKey == "" {
					err = fmt.Errorf("API key is missing, please save your API key in settings first")
				} else if cfg.ApiKey == "DRY_RUN" {
					// Simulate processing delay
					responseTimeMs = int64(1000 + (vIdx+mIdx)*250%1000) // 1.0s to 2.0s
					
					// Sleep with cancel check support
					select {
					case <-time.After(time.Duration(responseTimeMs) * time.Millisecond):
					case <-processCtx.Done():
						a.emitEvent("task:cancelled")
						return
					}

					resultText = fmt.Sprintf("[GIẢ LẬP DRY RUN - MÔ HÌNH %s]\nVăn bản %d\n---\nNội dung gốc: %s\n\nĐây là kết quả tóm tắt giả lập chất lượng cao từ mô hình %s.", 
						modelName, vIdx+1, truncateText(trimmedText, 60), modelName)
				} else {
					client := backend.NewGeminiClient(cfg.ApiKey)
					resultText, responseTimeMs, err = client.SummarizeText(modelName, cfg.SystemPrompt, trimmedText)
				}

				if err != nil {
					a.emitEvent("task:failed", map[string]interface{}{
						"modelIdx":     mIdx,
						"vbIdx":        vIdx,
						"errorMessage": err.Error(),
					})
				} else {
					errSave := backend.SaveResult(cfg.OutputDir, mIdx+1, vIdx+1, resultText, responseTimeMs)
					if errSave != nil {
						a.emitEvent("task:failed", map[string]interface{}{
							"modelIdx":     mIdx,
							"vbIdx":        vIdx,
							"errorMessage": fmt.Sprintf("failed to save result: %v", errSave),
						})
					} else {
						a.emitEvent("task:success", map[string]interface{}{
							"modelIdx":       mIdx,
							"vbIdx":          vIdx,
							"responseTimeMs": responseTimeMs,
						})
					}
				}
			}
		}

		// Save final mapping.json upon success
		_ = backend.SaveMapping(cfg.OutputDir, cfg.Models, texts)
	}()
}

// CancelBatchProcess cancels the currently active batch summary execution
func (a *App) CancelBatchProcess() {
	if a.cancelFunc != nil {
		a.cancelFunc()
	}
}
