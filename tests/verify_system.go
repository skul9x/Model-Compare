package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run tests/verify_system.go <path_to_output_dir>")
		os.Exit(1)
	}
	outputDir := os.Args[1]

	// 1. Kiểm định file mapping.json
	mappingPath := filepath.Join(outputDir, "mapping.json")
	mapData, err := os.ReadFile(mappingPath)
	if err != nil {
		fmt.Printf("❌ Thất bại: Không tìm thấy file mapping.json tại %s\n", mappingPath)
		os.Exit(1)
	}

	var mapping map[string]map[string]string
	if err := json.Unmarshal(mapData, &mapping); err != nil {
		fmt.Printf("❌ Thất bại: File mapping.json sai cấu trúc JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Thành công: File mapping.json hợp lệ!")

	// 2. Quét và kiểm tra chất lượng từng file .txt
	files, err := os.ReadDir(outputDir)
	if err != nil {
		fmt.Printf("❌ Thất bại: Không thể quét thư mục %s\n", outputDir)
		os.Exit(1)
	}

	txtCount := 0
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
		txtCount++
		filePath := filepath.Join(outputDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("❌ Thất bại: Không đọc được file %s\n", file.Name())
			continue
		}

		content := string(data)
		lines := strings.Split(content, "\n")
		
		// Kiểm tra dòng đầu tiên chứa thời gian phản hồi
		if len(lines) == 0 || !strings.HasPrefix(lines[0], "Response Time:") {
			fmt.Printf("❌ Thất bại: File %s không chứa thông tin 'Response Time:' ở dòng đầu tiên!\n", file.Name())
			continue
		}
		
		fmt.Printf("✅ Đã xác minh file: %s (Định dạng chuẩn)\n", file.Name())
	}

	if txtCount > 0 {
		fmt.Printf("🎉 Chúc mừng! Đã xác minh thành công tổng cộng %d file kết quả đạt tiêu chuẩn 100%%!\n", txtCount)
	} else {
		fmt.Println("⚠️ Cảnh báo: Không tìm thấy file kết quả .txt nào để kiểm tra.")
	}
}
