# Plan: Custom App Enhancements
Created: 2026-05-17T19:45:00+07:00
Status: 🟢 Completed

## Overview
Lên kế hoạch nâng cấp giao diện và tính năng cho ứng dụng Google Model Tester theo yêu cầu mới từ người dùng:
1. Thay thế khu vực nhập văn bản đơn lẻ (`#bulk-text-input`) bằng danh sách các khu vực nhập văn bản (`textarea`) riêng biệt và động, hỗ trợ thêm/xóa từng đoạn văn bản trực quan qua nút bấm thay vì phân cách bằng ký tự `---`.
2. **[MỚI]** Thêm nút **Sao chép (Copy)** trên mỗi thẻ nhập văn bản để copy toàn bộ nội dung của ô nhập đó một cách nhanh chóng.
3. **[MỚI]** Thêm khu vực **Ghi chú ký hiệu file (Legend)** giúp người dùng biết chính xác ký hiệu `model_1`, `model_2` trong file đầu ra ứng với tên thật của mô hình nào.
4. Nâng cấp và cải tiến Default System Prompt thành phiên bản chuyên nghiệp hơn, giúp tóm tắt cô đọng, giữ nguyên ngôn ngữ gốc và có định dạng gạch đầu dòng rõ ràng.
5. Đơn giản hóa phần thông tin Copyright ở Footer chỉ hiển thị duy nhất "Nguyễn Duy Trường 2026", loại bỏ mọi chi tiết khác liên quan đến IDE hay Antigravity.

## Tech Stack
- **Framework:** Wails v2 (Go Backend + HTML/JS/CSS Frontend)
- **Frontend Core:** Vanilla HTML5, CSS3 Custom Properties, Vanilla JavaScript (ES6)

## Phases

| Phase | Name | Status | Progress |
|-------|------|--------|----------|
| 01 | Cập nhật cấu trúc UI & CSS (HTML/CSS) | 🟢 Completed | 100% |
| 02 | Tích hợp Logic JS điều khiển danh sách Textareas, Legend & Prompt | 🟢 Completed | 100% |
| 03 | Kiểm thử, Xác minh & Bàn giao | 🟢 Completed | 100% |

## Quick Commands
- Bắt đầu Phase 1: `/code phase-01`
- Xem gợi ý tiếp theo: `/next`
- Lưu lại trạng thái: `/save-brain`
