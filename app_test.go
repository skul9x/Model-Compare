package main

import (
	"context"
	"encoding/json"
	"modeltester/backend"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestApp_StartBatchProcess_DryRun(t *testing.T) {
	// Create a temp output directory
	tempDir, err := os.MkdirTemp("", "modeltester_app_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp output dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	app := NewApp()
	app.isTest = true
	app.ctx = context.Background()

	// Dry Run Configuration
	cfg := backend.AppConfig{
		ApiKey:       "DRY_RUN",
		Models:       []string{"mock-gemini-1.5-flash", "mock-gemini-1.5-pro"},
		SystemPrompt: "Hãy tóm tắt...",
		OutputDir:    tempDir,
	}

	texts := []string{
		"Văn bản số một cần tóm tắt.",
		"Văn bản số hai cũng cần tóm tắt.",
	}

	// We start the batch process. Since it runs in a goroutine, we need to wait for its completion.
	// Wails normally calls a.cancelFunc on complete and emits events.
	// Since there is no real Wails runtime, it won't panic, it will just log/warn or ignore.
	app.StartBatchProcess(cfg, texts)

	// Since Dry Run has 1000ms+ delay per text/model pair, let's wait for them to finish.
	// Wait, we can sleep or check the files being created!
	// There are 2 texts * 2 models = 4 tasks.
	// Expected files:
	// - model_1-vb_1.txt, model_2-vb_1.txt, model_1-vb_2.txt, model_2-vb_2.txt
	// - mapping.json
	
	t.Log("Waiting for dry run tasks to write output files...")
	
	// Set a timeout of 10 seconds for verification
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	expectedFiles := []string{
		filepath.Join(tempDir, "model_1-vb_1.txt"),
		filepath.Join(tempDir, "model_1-vb_2.txt"),
		filepath.Join(tempDir, "model_2-vb_1.txt"),
		filepath.Join(tempDir, "model_2-vb_2.txt"),
		filepath.Join(tempDir, "mapping.json"),
	}

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for batch process to write all expected files")
		case <-ticker.C:
			allCreated := true
			for _, file := range expectedFiles {
				if _, err := os.Stat(file); os.IsNotExist(err) {
					allCreated = false
					break
				}
			}
			if allCreated {
				// Also check if the model_{time}.txt is created
				timeFiles, _ := filepath.Glob(filepath.Join(tempDir, "model_[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9][0-9][0-9].txt"))
				if len(timeFiles) > 0 {
					t.Log("All files and model_{time}.txt created successfully!")
					goto VERIFY
				}
			}
		}
	}

VERIFY:
	// Verify that the contents of one of the files matches dry run simulation format
	data, err := os.ReadFile(expectedFiles[0])
	if err != nil {
		t.Fatalf("Failed to read model result file: %v", err)
	}

	contentStr := string(data)
	if !contains(contentStr, "[GIẢ LẬP DRY RUN - MÔ HÌNH mock-gemini-1.5-flash]") {
		t.Errorf("Result file doesn't contain expected Dry Run header: %s", contentStr)
	}

	// Verify mapping file
	mappingData, err := os.ReadFile(filepath.Join(tempDir, "mapping.json"))
	if err != nil {
		t.Fatalf("Failed to read mapping.json: %v", err)
	}

	var mapping map[string]interface{}
	if err := json.Unmarshal(mappingData, &mapping); err != nil {
		t.Fatalf("Failed to parse mapping.json: %v", err)
	}

	modelsMap := mapping["models"].(map[string]interface{})
	if modelsMap["model_1"] != "mock-gemini-1.5-flash" || modelsMap["model_2"] != "mock-gemini-1.5-pro" {
		t.Errorf("Mapping models do not match expected: %v", modelsMap)
	}

	// Verify the model_{time}.txt file contents
	timeFiles, err := filepath.Glob(filepath.Join(tempDir, "model_[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9][0-9][0-9].txt"))
	if err != nil {
		t.Fatalf("Failed to glob time files: %v", err)
	}
	if len(timeFiles) == 0 {
		t.Error("Expected model_{time}.txt file to be created, but none found")
	} else {
		contentBytes, err := os.ReadFile(timeFiles[0])
		if err != nil {
			t.Fatalf("Failed to read time file: %v", err)
		}
		expectedContent := "model_1: mock-gemini-1.5-flash\nmodel_2: mock-gemini-1.5-pro\n"
		if string(contentBytes) != expectedContent {
			t.Errorf("Time file content is incorrect. Expected: %s, Got: %s", expectedContent, string(contentBytes))
		}
	}
}

func TestApp_CancelBatchProcess(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "modeltester_app_test_cancel_*")
	if err != nil {
		t.Fatalf("Failed to create temp output dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	app := NewApp()
	app.isTest = true
	app.ctx = context.Background()

	// Dry Run Configuration
	cfg := backend.AppConfig{
		ApiKey:       "DRY_RUN",
		Models:       []string{"mock-gemini-1.5-flash", "mock-gemini-1.5-pro"},
		SystemPrompt: "Hãy tóm tắt...",
		OutputDir:    tempDir,
	}

	// Supply lots of texts so it takes a long time
	texts := []string{
		"Văn bản 1", "Văn bản 2", "Văn bản 3", "Văn bản 4",
		"Văn bản 5", "Văn bản 6", "Văn bản 7", "Văn bản 8",
	}

	// Start the process
	app.StartBatchProcess(cfg, texts)

	// Sleep slightly to let the first model task start
	time.Sleep(100 * time.Millisecond)

	// Trigger Cancel
	t.Log("Cancelling batch process...")
	app.CancelBatchProcess()

	// Sleep to allow goroutine to terminate and cleanup
	time.Sleep(500 * time.Millisecond)

	// Since we cancelled, there should be fewer files than the total 16 files
	files, err := filepath.Glob(filepath.Join(tempDir, "model_*-vb_*.txt"))
	if err != nil {
		t.Fatalf("Failed to glob output files: %v", err)
	}

	t.Logf("Number of files written before cancel took effect: %d", len(files))
	if len(files) == 16 {
		t.Error("Expected progress to be cancelled, but all 16 files were written")
	}

	// Verify that the model_{time}.txt file is still created even when cancelled
	timeFiles, err := filepath.Glob(filepath.Join(tempDir, "model_[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9][0-9][0-9].txt"))
	if err != nil {
		t.Fatalf("Failed to glob time files upon cancel: %v", err)
	}
	if len(timeFiles) == 0 {
		t.Error("Expected model_{time}.txt file to be created upon cancellation, but none found")
	} else {
		contentBytes, err := os.ReadFile(timeFiles[0])
		if err != nil {
			t.Fatalf("Failed to read cancel time file: %v", err)
		}
		expectedContent := "model_1: mock-gemini-1.5-flash\nmodel_2: mock-gemini-1.5-pro\n"
		if string(contentBytes) != expectedContent {
			t.Errorf("Cancel time file content is incorrect. Expected: %s, Got: %s", expectedContent, string(contentBytes))
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && stringContains(s, substr)
}

func stringContains(s, substr string) bool {
	// Simple lookup
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestApp_StartBatchProcess_EmptyAndInvalidInputs(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "modeltester_app_test_inputs_*")
	if err != nil {
		t.Fatalf("Failed to create temp output dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	app := NewApp()
	app.isTest = true
	app.ctx = context.Background()

	// 1. Test case 1: Empty/whitespace texts
	cfg1 := backend.AppConfig{
		ApiKey:       "DRY_RUN",
		Models:       []string{"mock-model"},
		SystemPrompt: "Hãy tóm tắt...",
		OutputDir:    tempDir,
	}
	// Gửi cả văn bản hợp lệ và văn bản trống/khoảng trắng
	texts1 := []string{
		"Văn bản hợp lệ",
		"",
		"   ",
	}

	app.StartBatchProcess(cfg1, texts1)

	// Đợi cho task hoàn tất. Vì có 3 văn bản, văn bản 2 và 3 sẽ fail ngay lập tức, văn bản 1 mất 1 giây.
	// Chúng ta chỉ cần chờ file output được sinh ra.
	t.Log("Waiting for task output to verify empty texts handling...")
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// expectedFiles: chỉ model_1-vb_1.txt được tạo ra vì vb_2 và vb_3 trống nên sẽ bị lỗi (không lưu kết quả thành công)
	expectedSuccessFile := filepath.Join(tempDir, "model_1-vb_1.txt")
	unexpectedFile2 := filepath.Join(tempDir, "model_1-vb_2.txt")
	unexpectedFile3 := filepath.Join(tempDir, "model_1-vb_3.txt")

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for successful task to write output")
		case <-ticker.C:
			if _, err := os.Stat(expectedSuccessFile); err == nil {
				goto VERIFY_EMPTY
			}
		}
	}

VERIFY_EMPTY:
	// Xác minh vb_2 và vb_3 không tạo ra file output .txt thành công
	if _, err := os.Stat(unexpectedFile2); err == nil {
		t.Error("Expected no output file for empty text (vb_2), but it was created")
	}
	if _, err := os.Stat(unexpectedFile3); err == nil {
		t.Error("Expected no output file for whitespace text (vb_3), but it was created")
	}

	// 2. Test case 2: Empty API Key
	cfg2 := backend.AppConfig{
		ApiKey:       "",
		Models:       []string{"mock-model"},
		SystemPrompt: "Hãy tóm tắt...",
		OutputDir:    tempDir,
	}
	texts2 := []string{"Văn bản kiểm tra API Key rỗng"}

	// Tạo một thư mục output trống mới để kiểm tra dễ dàng
	tempDir2, err := os.MkdirTemp("", "modeltester_app_test_apikey_*")
	if err != nil {
		t.Fatalf("Failed to create temp output dir: %v", err)
	}
	defer os.RemoveAll(tempDir2)
	cfg2.OutputDir = tempDir2

	app.StartBatchProcess(cfg2, texts2)

	// Vì API Key rỗng, tác vụ sẽ fail ngay lập tức và không có file output nào được ghi.
	time.Sleep(500 * time.Millisecond)

	expectedFailedFile := filepath.Join(tempDir2, "model_1-vb_1.txt")
	if _, err := os.Stat(expectedFailedFile); err == nil {
		t.Error("Expected no output file for empty API key, but it was created")
	}
}

func TestSimulationForVerification(t *testing.T) {
	outputDir := "./test_output"
	// Đảm bảo dọn dẹp thư mục cũ trước khi chạy
	os.RemoveAll(outputDir)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test_output dir: %v", err)
	}

	app := NewApp()
	app.isTest = true
	app.ctx = context.Background()

	// Dry Run Configuration với 3 models
	cfg := backend.AppConfig{
		ApiKey:       "DRY_RUN",
		Models:       []string{"gemini-1.5-flash", "gemini-1.5-pro", "gemini-fake-model"},
		SystemPrompt: "Hãy tóm tắt văn bản sau thành danh sách...",
		OutputDir:    outputDir,
	}

	// 3 Đoạn văn bản (trong đó có 1 đoạn trống ở giữa để kiểm tra xử lý lỗi biên)
	texts := []string{
		"Văn bản thực tế thứ nhất có nội dung khá dài để kiểm nghiệm khả năng tóm tắt của các mô hình.",
		"   ", // Đoạn trống/khoảng trắng
		"Văn bản thực tế thứ ba ngắn gọn hơn nhưng vẫn đầy đủ cấu trúc.",
	}

	t.Log("Starting simulation for verification...")
	app.StartBatchProcess(cfg, texts)

	// Có 3 models * 3 texts = 9 tasks. 
	// Trong đó, các task ứng với văn bản 2 (chỉ khoảng trắng) sẽ thất bại ngay lập tức ở backend (3 tasks failed).
	// Các task ứng với văn bản 1 và 3 sẽ thành công và ghi file (6 tasks success).
	
	expectedSuccessFiles := []string{
		filepath.Join(outputDir, "model_1-vb_1.txt"),
		filepath.Join(outputDir, "model_2-vb_1.txt"),
		filepath.Join(outputDir, "model_3-vb_1.txt"),
		filepath.Join(outputDir, "model_1-vb_3.txt"),
		filepath.Join(outputDir, "model_2-vb_3.txt"),
		filepath.Join(outputDir, "model_3-vb_3.txt"),
		filepath.Join(outputDir, "mapping.json"),
	}

	timeout := time.After(12 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for simulation files to be written")
		case <-ticker.C:
			allCreated := true
			for _, file := range expectedSuccessFiles {
				if _, err := os.Stat(file); os.IsNotExist(err) {
					allCreated = false
					break
				}
			}
			if allCreated {
				t.Log("Simulation completed! All expected files are written successfully.")
				return
			}
		}
	}
}
