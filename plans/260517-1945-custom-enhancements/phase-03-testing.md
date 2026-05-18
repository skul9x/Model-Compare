# Phase 03: Kiểm thử, Xác minh & Bàn giao
Status: 🟢 Completed
Dependencies: [Phase 02: Tích hợp Logic JS điều khiển danh sách Textareas, Legend & Prompt](file:///home/skul9x/Desktop/Test_code/Model-Compare/plans/260517-1945-custom-enhancements/phase-02-integration.md)

## Objective
Chạy thử ứng dụng bằng chế độ Dry Run (Sử dụng API Key là `DRY_RUN`) và API Key thật của Gemini để kiểm tra toàn bộ luồng hoạt động. Đảm bảo dữ liệu từ các textarea riêng biệt được truyền xuống Go Backend chính xác, ghi file kết quả và mapping.json đúng chuẩn, xác minh nút Copy hoạt động đúng, Legend hiển thị đồng bộ và kiểm tra tính thẩm mỹ của giao diện mới.

## Requirements
### Functional
- [x] Chế độ Dry Run hoạt động hoàn hảo với nhiều ô nhập văn bản động mà không gặp lỗi.
- [x] Tính năng sao chép nội dung của từng ô nhập liệu vào bộ nhớ đệm (Clipboard) hoạt động trơn tru.
- [x] Khu vực Legend (Ghi chú ký hiệu file) hiển thị đúng tên mô hình thực tế ứng với ký hiệu `model_1`, `model_2`...
- [x] Dữ liệu kết quả đầu ra lưu đúng thư mục đã chọn dưới dạng các file `model_X-vb_Y.txt`.
- [x] File `mapping.json` lưu đúng thông tin ánh xạ tương ứng với từng văn bản động: `vb_1`, `vb_2`...
- [x] Quá trình hủy tiến trình (Cancel) hoạt động bình thường trên danh sách ô nhập mới.

### Non-Functional (Verification)
- [x] Xác minh bản quyền hiển thị đơn giản ở chân trang đúng yêu cầu.
- [x] Đảm bảo hiệu năng tải trang mượt mà khi người dùng tạo nhiều ô nhập (lên tới 10-20 ô nhập cùng lúc).

## Implementation Steps
1. [x] **Khởi động ứng dụng bằng Wails Dev**:
   - Sử dụng lệnh `wails dev` từ thư mục gốc để biên dịch và chạy ứng dụng ở chế độ Live Development.
2. [x] **Thực hiện bài kiểm thử thủ công (Manual Testing)**:
   - **Test Case 1**: Kiểm tra mặc định có 1 ô nhập văn bản. Nút Xóa bị ẩn. Nút Sao chép (Copy) hiển thị đầy đủ.
   - **Test Case 2**: Click "Thêm đoạn văn bản" 3 lần để tạo tổng cộng 4 ô nhập. Con trỏ tự động chuyển đến ô mới nhất. Nút Xóa xuất hiện trên tất cả các ô.
   - **Test Case 3**: Nhập nội dung khác nhau vào cả 4 ô. Click thử nút "Sao chép" ở ô số 3, dán thử vào trình duyệt hoặc ứng dụng khác để xác minh đã sao chép đúng. Kiểm tra xem icon có đổi thành dấu tích xanh lá mượt mà trong 1.5 giây hay không.
   - **Test Case 4**: Nhìn vào card quản lý Models trong Sidebar, kiểm tra xem Legend ghi chú có hiển thị đúng khớp với danh sách models phía trên (Ví dụ: `model_1: gemini-1.5-flash`, `model_2: gemini-1.5-pro`). Thử thêm 1 model mới và xóa 1 model cũ xem Legend có tự động cập nhật đồng bộ hay không.
   - **Test Case 5**: Xóa ô nhập thứ 2 và thứ 3. Xác minh danh sách Preview tự cập nhật lại chỉ còn 2 văn bản và các nhãn đánh số thứ tự nhảy lại chính xác (vb_1, vb_2).
   - **Test Case 6**: Chạy tóm tắt giả lập bằng cách nhập API Key: `DRY_RUN`. Chọn một thư mục đầu ra bất kỳ. Click "Bắt đầu tóm tắt".
   - **Test Case 7**: Kiểm tra ma trận tiến trình chạy tuần tự từng tác vụ đến khi hoàn tất 100%.
3. [x] **Xác thực cấu trúc File đầu ra**:
   - Truy cập vào thư mục đầu ra đã chọn để kiểm tra các file lưu trữ:
     - Các file kết quả tóm tắt: `model_1-vb_1.txt`, `model_2-vb_1.txt`...
     - File ánh xạ: `mapping.json` chứa đúng nội dung đã nhập ở các ô tương ứng.

## Files to Create/Modify
- `tests/verify_phase3.js` - Tạo script kiểm thử tự động kiểm tra bản quyền chân trang, cấu trúc sinh mapping.json, đặt tên file output, khả năng mở rộng tối đa (10-20 ô nhập), sự hiện diện của API Copy Clipboard và chức năng Hủy (Cancel).

## Test Criteria
- [x] Tất cả các Test Cases đều vượt qua (Pass).
- [x] File `mapping.json` hiển thị đúng cấu trúc JSON ánh xạ tên mô hình và văn bản nhập liệu thực tế.
- [x] Không có lỗi runtime trên cả Go Console lẫn Web Browser Console (Wails Inspector).

---
Bàn giao ứng dụng hoàn thiện và sẵn sàng cho việc đóng gói sản phẩm.
