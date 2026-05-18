# Phase 01: Setup Environment & Wails Template
Status: ✅ Completed
Dependencies: None

## Objective
Khởi tạo dự án Go + Wails v2 với template Vanilla HTML/JS, dọn dẹp các file mặc định không dùng tới, cấu hình go.mod và kiểm tra khả năng chạy dev.

## Requirements
### Functional
- [x] Khởi tạo thành công cấu trúc dự án Wails v2 trong thư mục hiện tại.
- [x] Dọn dẹp giao diện mặc định (HTML, CSS, JS mặc định của Wails) để chuẩn bị cho giao diện custom.
- [x] Cấu hình chạy ứng dụng ở chế độ dev (`wails dev`) thành công.

### Non-Functional
- [x] Đảm bảo không có lỗi biên dịch Go hoặc NPM cài đặt thư viện.

## Implementation Steps
1. [x] Mở terminal tại thư mục dự án và khởi tạo dự án Wails:
   ```bash
   wails init -n modeltester -t vanilla
   ```
   *Lưu ý: Vì Wails v2 init sẽ tạo thư mục con `modeltester`, chúng ta sẽ chuyển các file ra thư mục gốc `d:\skul9x\Model-Tester` hoặc làm việc trong thư mục con này.*
2. [x] Di chuyển các file từ thư mục `modeltester` ra ngoài thư mục gốc để thư mục gốc chứa trực tiếp dự án (giúp đường dẫn workspace chuẩn nhất).
3. [x] Xóa các file frontend không cần thiết trong `frontend/src/` (như logo mặc định, các logic demo).
4. [x] Chạy lệnh chạy thử ở chế độ dev để kiểm tra:
   ```bash
   wails dev
   ```

## Files to Create/Modify
- `wails.json` - File cấu hình Wails (chỉnh sửa tiêu đề app, kích thước cửa sổ mặc định).
- `frontend/index.html` - File HTML giao diện gốc.
- `frontend/src/style.css` - File CSS giao diện gốc (sẽ viết lại toàn bộ).
- `frontend/src/main.js` - File JS logic giao diện (sẽ viết lại toàn bộ).

## Phase Test File: `tests/verify_setup.go`
Mỗi tác vụ của Phase 01 sẽ được tự động kiểm tra bằng cách chạy file test này:
```go
package main

import (
	"fmt"
	"os"
)

func main() {
	files := []string{
		"wails.json",
		"go.mod",
		"main.go",
		"app.go",
		"frontend/index.html",
	}
	allExist := true
	for _, f := range files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			fmt.Printf("❌ Thiếu file quan trọng: %s\n", f)
			allExist = false
		} else {
			fmt.Printf("✅ Đã tìm thấy: %s\n", f)
		}
	}
	if allExist {
		fmt.Println("🚀 Môi trường Wails Phase 01 đã sẵn sàng hoàn hảo!")
	} else {
		fmt.Println("⚠️ Vui lòng hoàn thành các tác vụ còn thiếu.")
		os.Exit(1)
	}
}
```
*Cách chạy test:* `go run tests/verify_setup.go`

## Test Criteria
- [x] Chạy file `tests/verify_setup.go` trả về kết quả thành công `🚀 Môi trường Wails Phase 01 đã sẵn sàng hoàn hảo!`.
- [x] Lệnh `wails dev` biên dịch thành công mà không có lỗi.
- [x] Cửa sổ ứng dụng hiển thị lên với trang trắng hoặc tiêu đề ứng dụng mong muốn.

---
Next Phase: [Phase 02: Go Backend Core Logic (API & Files)](file:///d:/skul9x/Model-Tester/plans/260517-1745-google-model-tester/phase-02-backend.md)
