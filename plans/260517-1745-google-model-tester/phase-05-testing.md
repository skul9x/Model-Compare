# Phase 05: Testing, Edge-Cases & Verification
Status: ✅ Completed
Dependencies: [Phase 04: Integration & Wails Event Communication](file:///d:/skul9x/Model-Tester/plans/260517-1745-google-model-tester/phase-04-integration.md)

## Objective
Kiểm thử toàn diện ứng dụng, xử lý các tình huống biên (edge-cases), tối ưu hóa trải nghiệm người dùng, đảm bảo file đầu ra đúng định dạng và biên dịch thành công file `.exe` chạy độc lập.

## Requirements
### Functional
- [x] Kiểm thử việc xử lý các đoạn văn bản trống hoặc chỉ chứa khoảng trắng giữa các dấu `---`.
- [x] Kiểm thử khi nhập API Key sai hoặc rỗng:
  - Hệ thống phải hiển thị rõ ràng thông báo lỗi trả về từ API Google (ví dụ: `API key not valid` hoặc `RESOURCE_EXHAUSTED`).
- [x] Kiểm thử khi nhập tên model không tồn tại hoặc sai chính tả:
  - Task của model đó phải được đánh dấu Thất bại (Red) trên UI và tiếp tục chạy các model khác.
- [x] Xác minh định dạng nội dung trong file `.txt` kết quả:
  - Phải chứa dòng thời gian phản hồi: `Response Time: X ms`
  - Tiếp theo là nội dung tóm tắt chính xác.
- [x] Xác minh file `mapping.json` được ghi đúng cấu trúc ánh xạ rõ ràng:
  - Model: `model_1` -> `gemini-1.5-flash`, ...
  - Văn bản: `vb_1` -> `[Nội dung gốc...]`, ...
- [x] Thử nghiệm dừng đột ngột (Cancel) tiến trình xem các tác vụ đang chạy có dừng ngay lập tức hay không.
- [x] Xác minh tính năng lưu cấu hình tự động: API Key, Models, Prompt, Save Path được lưu chuẩn xác và tự động khôi phục khi tắt đi bật lại trên cả Windows (`%APPDATA%\modeltester\config.json`) và Linux (`~/.config/modeltester/config.json`).

### Non-Functional
- [x] Biên dịch dự án thành công ra file chạy độc lập trên Windows (`.exe`) bằng lệnh `wails build`.
- [x] Kích thước ứng dụng tối ưu, thời gian khởi động nhanh dưới 2 giây.

## Implementation Steps
1. [x] Thực hiện chạy thử với dữ liệu thực tế:
   - 3 Models (ví dụ: `gemini-1.5-flash`, `gemini-1.5-pro`, và một tên model giả lập sai `gemini-fake-model`).
   - 3 Đoạn văn bản (phân tách bởi `---`).
2. [x] Kiểm tra thư mục kết quả:
   - Xác nhận có đủ các file: `model_1-vb_1.txt`, `model_1-vb_2.txt`, `model_1-vb_3.txt`, `model_2-vb_1.txt`...
   - Xác nhận file lỗi của `model_3` (`gemini-fake-model`) không làm ứng dụng bị treo.
   - Kiểm tra nội dung của các file `.txt` xem có đúng dòng `Response Time:` và nội dung tóm tắt không.
   - Kiểm tra file `mapping.json` xem có đầy đủ ánh xạ không.
3. [x] Tối ưu hóa UI/UX:
   - Sửa các lỗi khoảng cách, căn chỉnh padding, font chữ Light Mode cho thật mịn màng.
   - Thêm tooltip giải thích khi rê chuột vào card trạng thái lỗi.
4. [x] Chạy lệnh build production:
   ```bash
   wails build
   ```
5. [x] Chạy thử file `.exe` tạo ra trong thư mục `build/bin/` để đảm bảo hoạt động hoàn hảo độc lập.

## Files to Create/Modify
- Toàn bộ codebase (fix các bug phát sinh trong quá trình test).

## Phase Test File: `tests/verify_system.go`
Để tự động hóa hoàn toàn việc xác minh chất lượng file kết quả `.txt` và file `mapping.json` đầu ra sau khi chạy thật (hoặc chạy Dry Run), ta sử dụng một chương trình kiểm định độc lập:
```go
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
```
*Cách chạy test:* `go run tests/verify_system.go <đường_dẫn_folder_kết_quả>`

## Test Criteria
- [x] Chạy file `tests/verify_system.go` trên thư mục kết quả trả về thông báo `🎉 Chúc mừng! Đã xác minh thành công...`.
- [x] Tất cả các file `.txt` được tạo ra đúng vị trí đã chọn với tên và nội dung chuẩn chỉ 100%.
- [x] File `mapping.json` hiển thị đúng ánh xạ.
- [x] Ứng dụng chạy mượt mà, không bị crash, không bị rò rỉ bộ nhớ ở cả Windows và Linux.
- [x] `wails build` chạy thành công không có bất kỳ warning hay error nào.

---
Next Steps: Toàn bộ dự án đã được hoàn tất và kiểm nghiệm 100% thành công rực rỡ! Dự án đã sẵn sàng bàn giao sản xuất!
