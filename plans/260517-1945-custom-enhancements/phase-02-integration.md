# Phase 02: Tích hợp Logic JS điều khiển danh sách Textareas, Legend & Prompt
Status: 🟢 Completed
Dependencies: [Phase 01: Cập nhật cấu trúc UI & CSS (HTML/CSS)](file:///home/skul9x/Desktop/Test_code/Model-Compare/plans/260517-1945-custom-enhancements/phase-01-frontend.md)

## Objective
Lập trình logic JavaScript để tạo, cập nhật, xóa, và sao chép các ô nhập liệu động trong mảng `state.texts`. Đồng thời hiển thị động ghi chú ánh xạ Models (`model_1`, `model_2`...) mỗi khi danh sách Models thay đổi, nâng cấp Default Prompt và truyền dữ liệu chuẩn xác xuống Go Backend.

## Requirements
### Functional
- [x] Khởi tạo `state.texts` với ít nhất một phần tử trống `[""]` khi ứng dụng khởi chạy lần đầu.
- [x] Hàm `renderTextInputList()` hiển thị danh sách các ô nhập liệu tương ứng với `state.texts`.
- [x] Cập nhật trực tiếp nội dung khi người dùng gõ vào từng textarea riêng lẻ vào `state.texts[index]`.
- [x] **[MỚI]** Lắng nghe sự kiện click trên nút Sao chép (`.btn-copy-text`) của từng ô để ghi nội dung vào Clipboard bằng `navigator.clipboard.writeText` và hiển thị hiệu ứng tích xanh (✓) phản hồi nhanh.
- [x] Lắng nghe sự kiện click trên nút Xóa để loại bỏ đoạn văn bản khỏi mảng `state.texts` và cập nhật lại giao diện.
- [x] Lắng nghe sự kiện click trên nút Thêm để chèn thêm `""` vào `state.texts`, vẽ lại danh sách và tự động `.focus()` vào ô nhập mới tạo.
- [x] Tự động ẩn nút Xóa nếu danh sách chỉ còn lại duy nhất một ô nhập liệu.
- [x] **[MỚI]** Xây dựng hàm `renderModelLegend()` để tạo danh sách ghi chú ánh xạ các mô hình đã lưu: `model_1` ứng với tên thật nào, `model_2` ứng với tên thật nào. Kích hoạt hàm này mỗi khi danh sách models được render (`renderModels()`).
- [x] Thay đổi `OLD_DEFAULT_PROMPT` cũ thành `NEW_DEFAULT_PROMPT` chuyên nghiệp. Khi tải config từ backend, nếu phát hiện prompt trống hoặc bằng giá trị mặc định cũ, hãy ghi đè bằng prompt mới và tự động lưu cấu hình.

## Implementation Steps
1. [x] **Khai báo DOM Elements mới**:
   - Khai báo `#text-inputs-container` và `#add-text-btn`.
   - Khai báo `#model-legend-container` và `#model-legend-list`.
2. [x] **Xây dựng hàm xử lý danh sách nhập liệu**:
   - `renderTextInputList()`: tạo cấu trúc HTML cho từng textarea kèm theo nút sao chép, nút xóa, nhãn thứ tự.
   - Thêm bộ lắng nghe sự kiện `input` trên từng textarea để cập nhật state đồng thời kích hoạt `updateTextsAndPreview()`.
   - **Xử lý copy**: Thêm listener cho `.btn-copy-text` để copy nội dung của index tương ứng, tạm thời đổi icon thành icon checkmark trong 1.5 giây để phản hồi người dùng.
3. [x] **Xây dựng hàm `renderModelLegend()`**:
   - Duyệt qua `state.models`, sinh ra mã HTML gạch đầu dòng dạng: `<li><strong>model_${index+1}:</strong> ${model}</li>`.
   - Ghi vào `#model-legend-list`. Nếu không có model nào, ẩn cả khung Legend.
   - Gọi `renderModelLegend()` ở cuối hàm `renderModels()`.
4. [x] **Cập nhật hàm bắt đầu tóm tắt (`startSummary`)**:
   - Lọc bỏ các văn bản rỗng (chỉ chứa khoảng trắng) trước khi gửi dữ liệu sang backend: `state.texts.map(t => t.trim()).filter(t => t.length > 0)`.
   - Cảnh báo nếu sau khi lọc, không có văn bản nào hợp lệ.
5. [x] **Nâng cấp Default Prompt**:
   - Định nghĩa hằng số `NEW_DEFAULT_PROMPT` trong `main.js`.
   - Cập nhật hàm `loadConfig()` để tự động nhận dạng prompt cũ hoặc rỗng và thay thế bằng prompt mới chất lượng cao, sau đó gọi `saveAppConfig()`.

## Files to Create/Modify
- `frontend/src/main.js` - Sửa đổi toàn bộ logic quản lý văn bản đầu vào, thiết lập các trình bắt sự kiện cho các ô nhập động, lập trình chức năng Copy Clipboard, render khung Legend của Models, nâng cấp prompt.

## Test Criteria
- [ ] Nhập liệu vào các ô khác nhau cập nhật đúng danh sách preview ở dưới.
- [ ] Click nút "Sao chép" ở ô văn bản nào thì copy đúng và chỉ copy nội dung của ô đó vào bộ nhớ đệm (Clipboard), icon tạm đổi thành tích xanh lá đẹp mắt.
- [ ] Click nút "Thêm đoạn văn bản" tạo ra ô mới và con trỏ chuột tự động nhấp nháy trong ô đó.
- [ ] Click nút "Xóa" loại bỏ đúng ô tương ứng, và nút Xóa biến mất khi chỉ còn lại 1 ô duy nhất.
- [ ] Khung Ghi chú Models hiển thị chính xác danh sách models hiện có dạng `model_1: gemini-1.5-flash`, `model_2: gemini-1.5-pro`... và tự động cập nhật khi thêm hoặc xóa mô hình.
- [ ] Khi reset app hoặc khởi động lần đầu, Default Prompt hiển thị phiên bản tiếng Việt nâng cấp dài, chuyên nghiệp hơn.

---
Next Phase: [Phase 03: Kiểm thử, Xác minh & Bàn giao](file:///home/skul9x/Desktop/Test_code/Model-Compare/plans/260517-1945-custom-enhancements/phase-03-testing.md)
