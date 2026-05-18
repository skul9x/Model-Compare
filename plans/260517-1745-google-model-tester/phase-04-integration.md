# Phase 04: Integration & Wails Event Communication
Status: ✅ Completed
Dependencies: [Phase 03: Frontend UI Development (Premium Light Mode)](file:///home/skul9x/Desktop/Test_code/Model-Compare/plans/260517-1745-google-model-tester/phase-03-frontend.md)

## Objective
Kết nối giao diện Frontend với logic Backend bằng Go thông qua các hàm Wails Binding và cơ chế truyền tin thời gian thực Wails Events.

## Requirements
### Functional
- [x] Giao diện gọi hàm `SelectDirectory()` của Go và hiển thị chính xác đường dẫn được chọn lên UI.
- [x] Tải cấu hình khi khởi chạy: Khi ứng dụng khởi động, Frontend gọi hàm `LoadConfig()` của Go để lấy dữ liệu config JSON đã lưu và tự động hiển thị lên giao diện (API Key, danh sách Models, Prompt, Save Path).
- [x] Tự động lưu cấu hình khi có thay đổi: Khi người dùng sửa API Key, Thêm/Xóa Model, sửa System Prompt, hoặc chọn folder lưu mới, Frontend gọi hàm `SaveConfig(cfg)` để cập nhật ngay lập tức xuống file hệ thống.
- [x] Khi click "Bắt đầu tóm tắt", Frontend sẽ:
  - Validate dữ liệu đầu vào (phải có API Key, ít nhất 1 model, ít nhất 1 văn bản, và đã chọn thư mục).
  - Tách mảng các đoạn văn bản.
  - Gửi dữ liệu xuống backend qua một hàm khởi tạo tiến trình chạy ngầm: `StartBatchProcess(config)`.
- [x] Backend chạy ngầm xử lý tuần tự (hoặc song song có giới hạn) các cặp `[Model x Văn bản]`:
  - Trước khi gọi API: Emit event `task:started` với thông tin `{modelIdx, vbIdx}`.
  - Nếu thành công:
    - Lưu file `model_X-vb_Y.txt` chứa thời gian phản hồi + kết quả.
    - Emit event `task:success` với `{modelIdx, vbIdx, responseTimeMs}`.
  - Nếu thất bại:
    - Ghi nhận lỗi vào log console trên giao diện.
    - Bỏ qua model lỗi đó, tiếp tục chạy các model khác.
    - Emit event `task:failed` với `{modelIdx, vbIdx, errorMessage}`.
  - Khi hoàn thành tất cả: Emit event `task:completed`.
- [x] Frontend lắng nghe các event trên bằng `runtime.EventsOn` để cập nhật trạng thái trực quan của các card trong bảng Grid tiến trình và ghi log vào console theo thời gian thực.
- [x] Tự động cập nhật Progress Bar theo tỷ lệ phần trăm các task đã hoàn thành (`(số task đã chạy / tổng số task) * 100`).

### Non-Functional
- [x] Tiến trình chạy ngầm trong Go sử dụng `context` của Wails để giao tiếp an toàn.
- [x] Tránh lag hoặc đơ giao diện bằng cách không chạy các tiến trình API đồng bộ trên thread chính của UI.

## Implementation Steps
1. [x] Cập nhật file `app.go` để nhận dữ liệu từ Frontend, khởi chạy một Goroutine riêng cho tiến trình batch summary.
2. [x] Viết logic truyền tin thời gian thực bằng `runtime.EventsEmit` từ backend xuống frontend.
3. [x] Trong Frontend `main.js`, đăng ký lắng nghe các sự kiện:
   - `runtime.EventsOn("task:started", ...)`
   - `runtime.EventsOn("task:success", ...)`
   - `runtime.EventsOn("task:failed", ...)`
   - `runtime.EventsOn("task:completed", ...)`
4. [x] Viết logic cập nhật giao diện tương ứng với mỗi event:
   - Thêm spinner xoay khi started.
   - Chuyển màu xanh lá và hiển thị số giây phản hồi khi success.
   - Chuyển màu đỏ và hiển thị icon cảnh báo (có tooltip hiển thị chi tiết lỗi) khi failed.
   - In các dòng log màu tương ứng vào console log box.
5. [x] Cấu hình nút "Hủy tiến trình" (Cancel) để cho phép dừng đột ngột nếu người dùng muốn dừng giữa chừng (sử dụng Go context Cancel).
6. [x] Tích hợp tính năng tải cấu hình tự động: Trong `main.js`, gọi hàm `LoadConfig()` khi ứng dụng khởi chạy thành công (sự kiện `wails.ready` hoặc `DOMContentLoaded`) để khôi phục cấu hình API Key, Models, Prompt, Save Path từ file `config.json`.
7. [x] Tích hợp tự động lưu cấu hình: Lắng nghe sự kiện thay đổi (`change` hoặc `input`) của các ô cấu hình trên UI để tự động gọi hàm `SaveConfig()` lưu lại dữ liệu ngay lập tức.

## Files to Create/Modify
- `app.go` - Cấu hình hàm `StartBatchProcess` và quản lý Goroutine/Context.
- `frontend/src/main.js` - Kết nối Wails bindings, lắng nghe Events, cập nhật UI động.

## Phase Test File: `tests/integration_input.txt` & Chế độ Giả lập (Dry Run)
Để kiểm thử tích hợp (Integration) toàn bộ luồng dữ liệu thời gian thực giữa Go Backend và JS Frontend mà không bị phụ thuộc vào API Key Gemini thực tế, ta thiết kế thêm **Chế độ Giả lập (Dry Run)**:
1. **File dữ liệu đầu vào mẫu `tests/integration_input.txt`**:
   ```txt
   Văn bản thứ nhất cần tóm tắt. Đây là phần test cho vb_1.
   ---
   Văn bản thứ hai có độ dài trung bình. Đây là phần test cho vb_2.
   ---
   Văn bản thứ ba rất ngắn. Đây là phần test cho vb_3.
   ```
2. **Logic Giả lập:** Nếu API Key nhập là `DRY_RUN`, backend Go sẽ tự động bỏ qua việc gọi API thật, giả lập quá trình xử lý mất từ 1-2 giây cho mỗi task, trả về kết quả giả lập ngẫu nhiên và sinh đầy đủ file kết quả cùng file mapping.

*Cách chạy test:*
- Chạy ứng dụng bằng `wails dev`.
- Chọn thư mục lưu trữ kết quả.
- Nhập API Key là `DRY_RUN`.
- Thêm 2-3 Model mẫu (ví dụ: `gemini-1.5-flash`, `gemini-1.5-pro`).
- Dán nội dung file `tests/integration_input.txt` vào TextArea.
- Bấm nút "Bắt đầu tóm tắt" và quan sát tiến trình, Progress Bar, bảng Grid và logs nhảy thời gian thực.
- Xác nhận các file `model_1-vb_1.txt`... được tạo ra tại folder lưu trữ.

## Test Criteria
- [x] Chạy giả lập với API Key `DRY_RUN` thành công, bảng tiến trình Grid cập nhật tuần tự các ô từ Gray (Chờ) -> Yellow (Đang chạy) -> Green (Thành công) cùng số giây phản hồi thực tế.
- [x] Progress Bar tăng dần chính xác từ 0% lên 100%.
- [x] Bấm nút "Browse" hiển thị đúng hộp thoại chọn thư mục gốc của hệ điều hành và cập nhật đường dẫn lên UI.
- [x] Ghi nhận đầy đủ file kết quả trong folder đích khi chạy xong đúng định dạng `model_X-vb_Y.txt` và file `mapping.json`.

---
## Integration & Verification Results
Chúng ta đã viết các bài kiểm thử tự động toàn diện và đạt kết quả mỹ mãn 100%:
1. **Backend Integration & Unit Tests (`go test -v ./...`):**
   - `TestApp_StartBatchProcess_DryRun`: Kiểm thử việc chạy bất đồng bộ với Dry Run thành công, tạo ra đầy đủ các file kết quả và mapping với cấu trúc hoàn chỉnh.
   - `TestApp_CancelBatchProcess`: Kiểm thử tính năng hủy giữa chừng thông qua Go context, xác nhận tiến trình dừng đột ngột và dọn dẹp tài nguyên chính xác.
   - Kết quả: **PASS** (100% thành công).

2. **Frontend Logic Tests (`node tests/verify_logic.js`):**
   - Kiểm thử tách chuỗi văn bản, thêm/xóa model.
   - Kết quả: **PASS** (100% thành công).

---
Next Phase: [Phase 05: Testing, Edge-Cases & Verification](file:///home/skul9x/Desktop/Test_code/Model-Compare/plans/260517-1745-google-model-tester/phase-05-testing.md)
