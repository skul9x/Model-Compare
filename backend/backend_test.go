package backend

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveConfigAndLoadConfig(t *testing.T) {
	// Ghi đè đường dẫn config sang file tạm để không làm ảnh hưởng config thực tế của người dùng
	tempFile, err := os.CreateTemp("", "modeltester_test_config_*.json")
	if err != nil {
		t.Fatalf("Không tạo được file config tạm: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	os.Remove(tempPath) // Xóa file để kiểm tra xem SaveConfig có tự động tạo thư mục/file không

	configPathOverride = tempPath
	defer func() {
		configPathOverride = ""
		os.Remove(tempPath)
	}()

	// Kiểm tra tính năng lưu/tải config đa nền tảng
	cfg := AppConfig{
		ApiKey:       "test-api-key-123",
		Models:       []string{"gemini-1.5-flash", "gemini-1.5-pro"},
		SystemPrompt: "Hãy tóm tắt đoạn văn bản...",
		OutputDir:    "./test_output",
	}

	err = SaveConfig(cfg)
	if err != nil {
		t.Fatalf("Không thể lưu cấu hình: %v", err)
	}

	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("Không thể đọc cấu hình: %v", err)
	}

	if loaded.ApiKey != cfg.ApiKey {
		t.Errorf("Mong đợi ApiKey '%s', nhận được '%s'", cfg.ApiKey, loaded.ApiKey)
	}
	if len(loaded.Models) != len(cfg.Models) || loaded.Models[0] != cfg.Models[0] {
		t.Errorf("Danh sách models không khớp")
	}
}

func TestSaveResultAndMapping(t *testing.T) {
	// Kiểm tra ghi kết quả tóm tắt và file mapping
	tempDir, err := os.MkdirTemp("", "modeltester_test")
	if err != nil {
		t.Fatalf("Không tạo được thư mục tạm: %v", err)
	}
	defer os.RemoveAll(tempDir)

	modelIdx := 1
	vbIdx := 1
	summaryContent := "Đây là nội dung tóm tắt mẫu."
	responseTimeMs := int64(1420)

	err = SaveResult(tempDir, modelIdx, vbIdx, summaryContent, responseTimeMs)
	if err != nil {
		t.Fatalf("Không thể lưu file kết quả: %v", err)
	}

	// 1. Kiểm tra file kết quả
	expectedFileName := filepath.Join(tempDir, "model_1-vb_1.txt")
	data, err := os.ReadFile(expectedFileName)
	if err != nil {
		t.Fatalf("Không tìm thấy file kết quả: %v", err)
	}

	contentStr := string(data)
	if !strings.Contains(contentStr, "Response Time: 1.42s (1420ms)") {
		t.Errorf("Nội dung file thiếu thông tin Response Time chuẩn")
	}
	if !strings.Contains(contentStr, summaryContent) {
		t.Errorf("Nội dung file thiếu nội dung tóm tắt thực tế")
	}

	// 2. Kiểm tra file mapping
	models := []string{"gemini-1.5-flash"}
	texts := []string{"Nội dung gốc văn bản 1"}
	err = SaveMapping(tempDir, models, texts)
	if err != nil {
		t.Fatalf("Không thể lưu file mapping: %v", err)
	}

	mappingFile := filepath.Join(tempDir, "mapping.json")
	mapData, err := os.ReadFile(mappingFile)
	if err != nil {
		t.Fatalf("Không tìm thấy file mapping: %v", err)
	}

	var mapping map[string]interface{}
	err = json.Unmarshal(mapData, &mapping)
	if err != nil {
		t.Fatalf("Lỗi parse JSON file mapping: %v", err)
	}

	modelsMap := mapping["models"].(map[string]interface{})
	if modelsMap["model_1"] != "gemini-1.5-flash" {
		t.Errorf("Ánh xạ models trong file mapping sai lệch")
	}
}

func TestSaveModelsTimeFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "modeltester_test_timefile")
	if err != nil {
		t.Fatalf("Không tạo được thư mục tạm: %v", err)
	}
	defer os.RemoveAll(tempDir)

	models := []string{"gemini-1.5-flash", "gemini-1.5-pro", "gemini-2.0-flash-exp"}
	filePath, err := SaveModelsTimeFile(tempDir, models)
	if err != nil {
		t.Fatalf("Không thể tạo file model_time: %v", err)
	}

	// Xác minh file tồn tại
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("File model_time không tồn tại: %s", filePath)
	}

	// Xác minh tên file khớp pattern model_YYYYMMDD-HHmmss.txt
	baseName := filepath.Base(filePath)
	if !strings.HasPrefix(baseName, "model_") || !strings.HasSuffix(baseName, ".txt") {
		t.Errorf("Tên file không đúng định dạng mong đợi: %s", baseName)
	}

	// Đọc nội dung và xác minh nội dung đúng định dạng
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Không thể đọc file model_time: %v", err)
	}

	content := string(data)
	expectedContent := "model_1: gemini-1.5-flash\nmodel_2: gemini-1.5-pro\nmodel_3: gemini-2.0-flash-exp\n"
	if content != expectedContent {
		t.Errorf("Nội dung file không đúng.\nMong đợi:\n%s\nNhận được:\n%s", expectedContent, content)
	}
}

func TestGeminiClient_SummarizeText(t *testing.T) {
	// Khởi tạo mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Kiểm tra query parameter key
		if r.URL.Query().Get("key") != "mock-api-key" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": {"code": 401, "message": "API key not valid", "status": "UNAUTHENTICATED"}}`))
			return
		}

		// Kiểm tra endpoint path có chứa model name và generateContent không
		if !strings.Contains(r.URL.Path, "gemini-1.5-flash:generateContent") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Parse request body
		var req GeminiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Kiểm tra system instructions
		if req.SystemInstruction == nil || len(req.SystemInstruction.Parts) == 0 || req.SystemInstruction.Parts[0].Text != "Tóm tắt ngắn gọn" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Kiểm tra text input
		if len(req.Contents) == 0 || len(req.Contents[0].Parts) == 0 || req.Contents[0].Parts[0].Text != "Đoạn văn cần tóm tắt" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Trả về JSON mock phản hồi thành công từ Gemini
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"candidates": [
				{
					"content": {
						"parts": [
							{
								"text": "Nội dung đã được tóm tắt thành công"
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

	client := NewGeminiClient("mock-api-key")
	client.BaseURL = server.URL // Trỏ URL của client vào mock server

	// Test trường hợp thành công
	summary, duration, err := client.SummarizeText("gemini-1.5-flash", "Tóm tắt ngắn gọn", "Đoạn văn cần tóm tắt")
	if err != nil {
		t.Fatalf("SummarizeText failed: %v", err)
	}

	if summary != "Nội dung đã được tóm tắt thành công" {
		t.Errorf("Mong đợi kết quả tóm tắt '%s', nhận được '%s'", "Nội dung đã được tóm tắt thành công", summary)
	}

	if duration < 0 {
		t.Errorf("Mong đợi thời gian phản hồi lớn hơn hoặc bằng 0, nhận được %d ms", duration)
	}

	// Test trường hợp API Key không hợp lệ (lỗi UNAUTHENTICATED)
	badClient := NewGeminiClient("wrong-key")
	badClient.BaseURL = server.URL
	_, _, err = badClient.SummarizeText("gemini-1.5-flash", "Tóm tắt ngắn gọn", "Đoạn văn cần tóm tắt")
	if err == nil {
		t.Errorf("Mong đợi lỗi khi truyền API key sai, nhưng không có lỗi")
	} else if !strings.Contains(err.Error(), "UNAUTHENTICATED") {
		t.Errorf("Mong đợi thông báo lỗi chứa 'UNAUTHENTICATED', nhận được: %v", err)
	}
}
