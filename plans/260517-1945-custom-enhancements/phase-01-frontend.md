# Phase 01: Cập nhật cấu trúc UI & CSS (HTML/CSS)
Status: 🟢 Completed
Dependencies: None

## Objective
Thay đổi giao diện người dùng bằng cách gỡ bỏ ô nhập văn bản hàng loạt và thay bằng hệ thống danh sách các ô nhập văn bản (`textarea`) độc lập, thêm nút Sao chép (Copy) và Xóa cho mỗi ô nhập. Bổ sung khu vực Ghi chú ký hiệu file (Legend) cho Models ở sidebar và làm sạch Footer bản quyền.

## Requirements
### Functional
- [x] Gỡ bỏ thẻ `<textarea id="bulk-text-input">` trong `frontend/index.html`.
- [x] Thêm thẻ chứa danh sách các ô nhập liệu: `<div id="text-inputs-container"></div>`.
- [x] Thêm nút thêm đoạn văn bản mới: `<button id="add-text-btn">`.
- [x] Thêm khung hiển thị ghi chú ký hiệu model: `<div id="model-legend-container">` bên trong thẻ quản lý Models của Sidebar.
- [x] Thay đổi nội dung footer bản quyền chỉ hiển thị duy nhất: `<p>Nguyễn Duy Trường 2026</p>`.

### Non-Functional (Aesthetics)
- [x] Thiết kế thẻ `.text-input-item` đồng nhất với hệ thống Glassmorphism có sẵn.
- [x] Tích hợp nhóm nút tác vụ `.text-input-actions` gồm nút Sao chép (Copy) và nút Xóa (Delete) trên góc phải tiêu đề của mỗi textarea.
- [x] Nút Sao chép sử dụng màu xanh ngọc dịu nhẹ, nút Xóa dùng màu đỏ dịu khi hover.
- [x] Khung `.model-legend-container` được thiết kế gọn gàng ở cuối card Models, sử dụng phông chữ nhỏ (`font-size: 0.8rem`), chữ nghiêng nhẹ và đường viền chia tách mờ nhạt để không làm rối mắt.
- [x] Thêm hiệu ứng slide-in mượt mà khi người dùng thêm ô nhập liệu mới.

## Implementation Steps
1. [x] **Chỉnh sửa `frontend/index.html`**:
   - Thay thế cấu trúc của `#bulk-text-input` bằng cấu trúc container danh sách động `#text-inputs-container` và nút bấm thêm mới `#add-text-btn`.
   - Thêm phần `#model-legend-container` ngay dưới danh sách `#model-list` trong card quản lý Models.
   - Sửa footer copyright thành `Nguyễn Duy Trường 2026`.
2. [x] **Bổ sung CSS vào `frontend/src/style.css`**:
   - Thêm lớp CSS cho `.text-inputs-container`, `.text-input-item`, `.text-input-header`, `.text-input-label`, `.text-input-actions`, `.btn-copy-text`, `.btn-delete-text`, `.add-text-action`.
   - Thêm lớp CSS cho `.model-legend-container`, `.legend-title`, `.legend-list`, `.legend-item` để hiển thị ghi chú model.
   - Bổ sung hiệu ứng chuyển động xuất hiện `slideIn` cho các ô nhập liệu mới.

## Files to Create/Modify
- `frontend/index.html` - Sửa đổi cấu trúc nhập liệu văn bản, thêm Legend container và sửa footer copyright.
- `frontend/src/style.css` - Bổ sung styling mới cho các nút Sao chép, Xóa, danh sách động, và khung ghi chú Legend.

## Test Criteria
- [x] Giao diện hiển thị đúng khung chứa và nút bấm "Thêm đoạn văn bản".
- [x] Mỗi thẻ nhập liệu hiển thị đầy đủ tiêu đề cùng 2 nút tác vụ Sao chép và Xóa ở góc phải.
- [x] Khung ghi chú Models (Legend) hiển thị đẹp mắt ở chân thẻ quản lý Models.
- [x] Footer hiển thị chính xác "Nguyễn Duy Trường 2026".

---
Next Phase: [Phase 02: Tích hợp Logic JS điều khiển danh sách Textareas, Legend & Prompt](file:///home/skul9x/Desktop/Test_code/Model-Compare/plans/260517-1945-custom-enhancements/phase-02-integration.md)
