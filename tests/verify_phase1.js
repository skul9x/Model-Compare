const fs = require('fs');
const path = require('path');
const assert = require('assert');

console.log('🧪 Bắt đầu kiểm thử tự động Phase 01: Cấu trúc UI & CSS...');

const htmlPath = path.join(__dirname, '../frontend/index.html');
const cssPath = path.join(__dirname, '../frontend/src/style.css');

// 1. Kiểm tra sự tồn tại của các file
assert.ok(fs.existsSync(htmlPath), 'File index.html phải tồn tại');
assert.ok(fs.existsSync(cssPath), 'File style.css phải tồn tại');
console.log('✅ Đã tìm thấy các file index.html và style.css');

const htmlContent = fs.readFileSync(htmlPath, 'utf8');
const cssContent = fs.readFileSync(cssPath, 'utf8');

// 2. Kiểm tra HTML
// - Gỡ bỏ bulk-text-input
assert.ok(!htmlContent.includes('id="bulk-text-input"'), 'Không được còn textarea id="bulk-text-input"');
console.log('✅ Đã loại bỏ bulk-text-input thành công');

// - Có text-inputs-container
assert.ok(htmlContent.includes('id="text-inputs-container"'), 'Phải có container id="text-inputs-container"');
console.log('✅ Đã thêm text-inputs-container thành công');

// - Có add-text-btn
assert.ok(htmlContent.includes('id="add-text-btn"'), 'Phải có nút id="add-text-btn"');
console.log('✅ Đã thêm add-text-btn thành công');

// - Có model-legend-container
assert.ok(htmlContent.includes('id="model-legend-container"'), 'Phải có container id="model-legend-container"');
console.log('✅ Đã thêm model-legend-container thành công');

// - Footer copyright chỉ chứa Nguyễn Duy Trường 2026
assert.ok(htmlContent.includes('<p>Nguyễn Duy Trường 2026</p>'), 'Footer copyright phải đúng "Nguyễn Duy Trường 2026"');
assert.ok(!htmlContent.includes('Pair programmed with Antigravity'), 'Không được chứa thông tin Antigravity trong footer');
console.log('✅ Footer copyright hiển thị chính xác và sạch sẽ');

// 3. Kiểm tra CSS
const requiredCssClasses = [
    '.text-inputs-container',
    '.text-input-item',
    '.text-input-header',
    '.text-input-label',
    '.text-input-actions',
    '.btn-copy-text',
    '.btn-delete-text',
    '.add-text-action',
    '.model-legend-container',
    '.legend-title',
    '.legend-list',
    '.legend-item'
];

requiredCssClasses.forEach(cls => {
    assert.ok(cssContent.includes(cls), `CSS phải chứa class ${cls}`);
});
console.log('✅ Đã kiểm tra đầy đủ các class CSS mới trong style.css');

console.log('\n🎉 TOÀN BỘ CÁC BÀI KIỂM THỬ PHASE 01 ĐÃ ĐẠT 100%! 🚀');
