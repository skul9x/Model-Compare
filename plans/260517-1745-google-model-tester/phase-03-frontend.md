# Phase 03: Frontend UI Development (Premium Light Mode)
Status: ✅ Completed
Dependencies: [Phase 02: Go Backend Core Logic (API & Files)](file:///d:/skul9x/Model-Tester/plans/260517-1745-google-model-tester/phase-02-backend.md)

## Objective
Xây dựng giao diện ứng dụng (Frontend) theo phong cách **Premium Light Mode** cực kỳ hiện đại, gọn gàng, tinh tế, sử dụng Vanilla HTML, CSS và JS.

## Requirements
### Functional
- [x] Nhập và ẩn/hiện API Key Gemini.
- [x] Danh sách Model:
  - Hiển thị danh sách model có sẵn.
  - Cho phép thêm model mới bằng ô nhập tên + nút "Add Model".
  - Cho phép xóa model khỏi danh sách.
- [x] Nhập văn bản hàng loạt:
  - Một TextArea lớn cho phép paste toàn bộ văn bản.
  - Tự động tách văn bản thành các đoạn `vb_1`, `vb_2`,... dựa trên ký tự phân tách `---` (đứng riêng trên 1 dòng).
  - Hiển thị trực quan danh sách các đoạn văn bản đã tách (`vb_1`, `vb_2`...) kèm theo ký tự mở đầu để người dùng preview trước khi chạy.
- [x] Nhập Custom System Prompt cho việc tóm tắt (Mặc định gợi ý: "Hãy tóm tắt đoạn văn bản sau đây thành một danh sách gạch đầu dòng được đánh số thứ tự 1, 2, 3... tập trung vào các ý chính quan trọng nhất.").
- [x] Chọn thư mục đầu ra qua nút "Browse" (gọi hàm `SelectDirectory` từ backend).
- [x] Nút "Bắt đầu tóm tắt" nổi bật với hiệu ứng hover mượt mà.
- [x] Panel tiến trình:
  - Một thanh Progress Bar tổng thể (tỷ lệ % hoàn thành).
  - Bảng trạng thái/Grid hiển thị danh sách tất cả các task `[Văn bản x Model]` dưới dạng các ô nhỏ (Card) có màu sắc thể hiện trạng thái: Chờ chạy (Gray), Đang chạy (Yellow - Spin), Thành công (Green + Response Time), Thất bại (Red + Tooltip lỗi).
  - Ô log console dạng văn bản hiển thị chi tiết tiến trình theo thời gian thực.

### Design & Aesthetics (Premium Light Mode)
- [x] **Color Palette:** 
  - Background chính: Trắng tinh khiết `#FFFFFF` và Xám siêu nhạt `#F8FAFC`.
  - Màu thương hiệu/Accents: Xanh ngọc thanh lịch `#0ea5e9` (Sky Blue) và Tím nhạt nhã nhặn `#8b5cf6`.
  - Văn bản: Xám đậm `#1e293b` (màu chính) và Xám nhạt `#64748b` (sub-title).
- [x] **Typography:** Sử dụng font chữ "Inter" hoặc "Outfit" từ Google Fonts cho cảm giác hiện đại, sạch sẽ.
- [x] **Premium Elements:** 
  - Glassmorphism nhẹ nhàng cho các card điều khiển (`backdrop-filter: blur(10px)` kết hợp viền trắng mờ).
  - Bo góc tinh tế (`border-radius: 12px` or `16px`).
  - Đổ bóng mềm mịn (`box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.05), 0 2px 4px -2px rgb(0 0 0 / 0.05)`).
  - Micro-animations: Hiệu ứng co giãn nhẹ khi di chuột qua nút (scale 1.02), danh sách model trượt nhẹ khi thêm/xóa.

## Implementation Steps
1. [x] Cài đặt Google Font "Inter" trong `frontend/index.html`.
2. [x] Thiết kế layout chính chia làm 2 cột chính trong `frontend/index.html`:
   - **Cột Trái (Cấu hình):** API Key, Model Management, System Prompt, Output Path.
   - **Cột Phải (Dữ liệu & Tiến trình):** Text Area nhập văn bản, Preview số lượng văn bản, Bảng Grid tiến trình trực quan, và Hộp Console Log.
3. [x] Viết CSS trong `frontend/src/style.css` định hình toàn bộ phong cách Light Mode cao cấp, responsive và micro-animations.
4. [x] Viết JS trong `frontend/src/main.js` xử lý:
   - Thêm/Xóa Model vào danh sách hiển thị trên UI.
   - Tự động lắng nghe sự thay đổi của TextArea văn bản để cập nhật Preview số lượng đoạn `vb` được chia.
   - Gọi các hàm binding Go để tải cấu hình khi khởi động và tự động cập nhật cấu hình mới mỗi khi người dùng thay đổi (API Key, Models, Prompt, Save Path) giúp đồng bộ đa nền tảng (Windows & Linux).

## Files to Create/Modify
- `frontend/index.html` - Giao diện HTML.
- `frontend/src/style.css` - CSS thiết kế giao diện Light Mode.
- `frontend/src/main.js` - Xử lý tương tác giao diện và quản lý state cục bộ.

## Phase Test File: `tests/verify_frontend.html`
Để kiểm duyệt trước giao diện (UI) và tính năng tách văn bản mà chưa cần chạy ứng dụng Wails đầy đủ, ta tạo một file kiểm thử UI độc lập. Anh chỉ cần mở file này trực tiếp bằng trình duyệt Chrome/Edge/Firefox để kiểm tra:
- **Tính năng tách văn bản:** Nhập dữ liệu có chứa `---` và xem nó đếm số đoạn chính xác.
- **Tính năng Models:** Thêm, xóa model trực quan xem danh sách có cập nhật mượt mà không.
- **Thiết kế Premium Light Mode:** Thử nghiệm xem màu sắc, độ bo góc, đổ bóng và độ mượt của các hover micro-animation.

*Cách chạy test:* Mở file `tests/verify_frontend.html` trên trình duyệt bằng cách kéo thả hoặc click đúp.

## Test Criteria
- [x] Mở file `tests/verify_frontend.html` trên trình duyệt hiển thị đúng chuẩn phong cách Light Mode cao cấp (màu xám nhạt, viền sky blue mịn màng, font chữ Inter sắc nét).
- [x] Giao diện hiển thị chuẩn, không bị vỡ khung hình khi co giãn kích thước màn hình trình duyệt.
- [x] Tính năng thêm/xóa Model hoạt động hoàn hảo trên UI mockup.
- [x] Paste các văn bản có chứa `---` tự động đếm và hiển thị đúng danh sách preview các đoạn văn bản.

---
Next Phase: [Phase 04: Integration & Wails Event Communication](file:///d:/skul9x/Model-Tester/plans/260517-1745-google-model-tester/phase-04-integration.md)
