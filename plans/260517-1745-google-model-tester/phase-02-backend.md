# Phase 02: Go Backend Core Logic (API & Files)
Status: ✅ Completed
Dependencies: [Phase 01: Setup Environment & Wails Template](file:///d:/skul9x/Model-Tester/plans/260517-1745-google-model-tester/phase-01-setup.md)

## Objective
Xây dựng logic Backend bằng Go bao gồm: 
1. Kết nối HTTP Direct tới Google Gemini API (để dễ dàng thay đổi bất kỳ model name nào).
2. Logic lưu file `.txt` kết quả theo chuẩn `model_X-vb_Y.txt` chứa thời gian phản hồi và nội dung tóm tắt.
3. Tạo file `mapping.json` lưu thông tin ánh xạ chi tiết giữa mã model/văn bản với tên thực tế.
4. Tích hợp bộ chọn thư mục native thông qua Wails Runtime.

## Requirements
### Functional
- [x] Gửi request HTTP POST tới `https://generativelanguage.googleapis.com/v1beta/models/{model}:generateContent?key={apiKey}`.
- [x] Tính toán thời gian phản hồi (Response Time) chính xác bằng Go (`time.Now()` trước và sau khi gọi API).
- [x] Ghi nhận log lỗi của model bị hỏng nhưng không làm crash app (trả về struct lỗi chi tiết).
- [x] Tạo folder lưu kết quả và ghi file `.txt` có định dạng:
  ```txt
  Response Time: 1.42s (1420ms)
  ---
  [NỘI DUNG TÓM TẮT]
  ```
- [x] Tạo và cập nhật file `mapping.json` chứa:
  - Danh sách model: `{"model_1": "gemini-1.5-flash", ...}`
  - Nội dung gốc của văn bản: `{"vb_1": "[Đoạn đầu văn bản gốc...]", ...}`
- [x] Cung cấp hàm `SelectDirectory()` trả về đường dẫn folder người dùng chọn.
- [x] Lưu cấu hình người dùng (API Key, danh sách các Google Models, System Prompt, Thư mục lưu đầu ra) đa nền tảng (Windows & Linux) bằng cách sử dụng thư mục cấu hình chuẩn của hệ điều hành thông qua hàm `os.UserConfigDir()`. Cấu hình sẽ được lưu vào file `config.json` nằm trong thư mục con `modeltester/`.

### Non-Functional
- [x] Code Go sạch, chia struct rõ ràng, xử lý lỗi (error handling) đầy đủ.
- [x] Sử dụng goroutine khi thực hiện chạy hàng loạt để tránh block UI.

## Implementation Steps
1. [x] Tạo struct `GeminiClient` trong Go để quản lý API Key và gửi request.
2. [x] Thiết lập struct Request/Response khớp với định dạng JSON của Gemini API.
3. [x] Viết hàm `SummarizeText(modelName string, systemPrompt string, text string) (string, int64, error)`:
   - Đo thời gian chạy.
   - Gọi API Gemini.
   - Parse kết quả trả về.
4. [x] Viết hàm `SaveResult(outputDir string, modelIdx int, vbIdx int, content string, responseTimeMs int64) error`.
5. [x] Viết hàm `SaveMapping(outputDir string, models []string, texts []string) error` để tạo file mapping lưu trữ thông tin thực tế.
6. [x] Tạo hàm Wails Binding `SelectDirectory()` sử dụng `wails runtime.OpenDirectoryDialog`.
7. [x] Tạo struct `AppConfig` và viết các hàm `LoadConfig() (AppConfig, error)` và `SaveConfig(cfg AppConfig) error`:
   - Sử dụng `os.UserConfigDir()` để lấy thư mục config chuẩn (Windows: `%APPDATA%\modeltester\config.json`, Linux: `~/.config/modeltester/config.json`).
   - Tự động dùng `os.MkdirAll(appConfigDir, 0700)` để tạo thư mục nếu chưa tồn tại.
   - Ghi/Đọc dữ liệu dưới định dạng JSON.

## Files to Create/Modify
- `backend/gemini.go` - Client kết nối Gemini API.
- `backend/files.go` - Logic lưu file kết quả và file mapping.
- `app.go` - File chính của Wails, binding các hàm Go cho Frontend gọi.

## Phase Test File: `backend/backend_test.go`
Mỗi logic nghiệp vụ của Backend Go được kiểm thử bằng unit test chuẩn chỉ:
```go
package backend

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveConfigAndLoadConfig(t *testing.T) {
	// Kiểm tra tính năng lưu/tải config đa nền tảng
	cfg := AppConfig{
		ApiKey:       "test-api-key-123",
		Models:       []string{"gemini-1.5-flash", "gemini-1.5-pro"},
		SystemPrompt: "Hãy tóm tắt đoạn văn bản...",
		OutputDir:    "./test_output",
	}

	err := SaveConfig(cfg)
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
```
*Cách chạy test:* Chạy trong terminal `go test ./backend -v`

## Test Criteria
- [x] Chạy thành công lệnh test backend `go test ./backend -v` mà không có bất kỳ lỗi nào.
- [x] Xác nhận file test tạo ra các file `.txt` và file `mapping.json` có cấu trúc hoàn chỉnh đúng yêu cầu của người dùng.

---
Next Phase: [Phase 03: Frontend UI Development (Premium Light Mode)](file:///d:/skul9x/Model-Tester/plans/260517-1745-google-model-tester/phase-03-frontend.md)
