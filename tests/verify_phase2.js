const assert = require('assert');

// 1. Core Logic implementations simulating frontend behaviour
const OLD_DEFAULT_PROMPT = 'Hãy tóm tắt đoạn văn bản sau đây thành một danh sách gạch đầu dòng được đánh số thứ tự 1, 2, 3... tập trung vào các ý chính quan trọng nhất.';
const NEW_DEFAULT_PROMPT = `Bạn là một chuyên gia phân tích và tóm tắt văn bản hiệu năng cao. Hãy tóm tắt đoạn văn bản được cung cấp đáp ứng các tiêu chuẩn sau:
1. Độ chính xác & Cô đọng: Trích xuất các ý chính, thông điệp quan trọng và số liệu cốt lõi một cách súc tích nhất, loại bỏ các chi tiết thừa thãi.
2. Định dạng rõ ràng: Trình bày kết quả dưới dạng danh sách gạch đầu dòng chuyên nghiệp.
3. Giữ nguyên ngôn ngữ gốc: Thực hiện tóm tắt bằng chính ngôn ngữ của văn bản đầu vào.
4. Tính trung thực: Không tự ý thêm thắt các thông tin, nhận định hoặc suy diễn ngoài văn bản gốc.`;

// Simulates loadConfig prompt upgrade
function upgradeSystemPrompt(loadedPrompt) {
    const rawPrompt = loadedPrompt ? loadedPrompt.trim() : '';
    let needsSave = false;
    let systemPrompt = '';
    
    if (rawPrompt === '' || rawPrompt === OLD_DEFAULT_PROMPT) {
        systemPrompt = NEW_DEFAULT_PROMPT;
        needsSave = true;
    } else {
        systemPrompt = loadedPrompt;
    }
    return { systemPrompt, needsSave };
}

// Simulates startSummary whitespace and empty text filtering
function filterValidTexts(texts) {
    return texts.map(t => t.trim()).filter(t => t.length > 0);
}

// Simulates Delete Button visibility rule (hide if length is 1)
function shouldHideDeleteButton(textsLength) {
    return textsLength === 1;
}

// Simulates model legend generation
function generateModelLegend(models) {
    if (models.length === 0) {
        return { hidden: true, items: [] };
    }
    const items = models.map((model, index) => {
        return {
            key: `model_${index + 1}`,
            value: model,
            html: `<span class="legend-key"><strong>model_${index + 1}:</strong></span> <span class="legend-val">${model}</span>`
        };
    });
    return { hidden: false, items };
}

// 2. Run Tests
function runTests() {
    console.log('🧪 Bắt đầu chạy các bài kiểm thử Phase 02 Integration...');

    // Test Case 1: Prompt Upgrade
    try {
        // Case 1.1: Empty prompt -> upgrades and saves
        const res1 = upgradeSystemPrompt('');
        assert.strictEqual(res1.systemPrompt, NEW_DEFAULT_PROMPT);
        assert.strictEqual(res1.needsSave, true);
        console.log('✅ Test Case 1.1: Nâng cấp prompt từ rỗng - ĐẠT');

        // Case 1.2: Old default prompt -> upgrades and saves
        const res2 = upgradeSystemPrompt(OLD_DEFAULT_PROMPT);
        assert.strictEqual(res2.systemPrompt, NEW_DEFAULT_PROMPT);
        assert.strictEqual(res2.needsSave, true);
        console.log('✅ Test Case 1.2: Nâng cấp prompt từ bản mặc định cũ - ĐẠT');

        // Case 1.3: User custom prompt -> does not upgrade, does not save
        const customPrompt = 'Tóm tắt ngắn gọn và vui nhộn nhé!';
        const res3 = upgradeSystemPrompt(customPrompt);
        assert.strictEqual(res3.systemPrompt, customPrompt);
        assert.strictEqual(res3.needsSave, false);
        console.log('✅ Test Case 1.3: Giữ nguyên prompt tùy chỉnh của người dùng - ĐẠT');
    } catch (e) {
        console.error('❌ Test Case 1: Thất bại:', e.stack);
        process.exit(1);
    }

    // Test Case 2: Whitespace and empty text filtering
    try {
        const inputTexts = [
            '  Đoạn văn bản số 1  ',
            '   ', // whitespace only
            'Đoạn văn bản số 2',
            '', // empty
            '\n\t\n' // whitespace only
        ];
        const validTexts = filterValidTexts(inputTexts);
        assert.strictEqual(validTexts.length, 2);
        assert.strictEqual(validTexts[0], 'Đoạn văn bản số 1');
        assert.strictEqual(validTexts[1], 'Đoạn văn bản số 2');
        console.log('✅ Test Case 2: Lọc văn bản rỗng và khoảng trắng chuẩn - ĐẠT');
    } catch (e) {
        console.error('❌ Test Case 2: Thất bại:', e.stack);
        process.exit(1);
    }

    // Test Case 3: Delete Button Visibility rule
    try {
        assert.strictEqual(shouldHideDeleteButton(1), true);
        assert.strictEqual(shouldHideDeleteButton(2), false);
        assert.strictEqual(shouldHideDeleteButton(5), false);
        console.log('✅ Test Case 3: Ẩn nút xóa khi chỉ có 1 ô nhập liệu - ĐẠT');
    } catch (e) {
        console.error('❌ Test Case 3: Thất bại:', e.stack);
        process.exit(1);
    }

    // Test Case 4: Model Legend generation
    try {
        // Case 4.1: Empty models list -> legend container hidden
        const legend1 = generateModelLegend([]);
        assert.strictEqual(legend1.hidden, true);
        assert.strictEqual(legend1.items.length, 0);
        console.log('✅ Test Case 4.1: Ẩn ghi chú Legend khi không có model - ĐẠT');

        // Case 4.2: Multiple models -> correctly mapped keys and HTML output
        const models = ['gemini-1.5-flash', 'gemini-1.5-pro', 'claude-3-opus'];
        const legend2 = generateModelLegend(models);
        assert.strictEqual(legend2.hidden, false);
        assert.strictEqual(legend2.items.length, 3);
        
        assert.strictEqual(legend2.items[0].key, 'model_1');
        assert.strictEqual(legend2.items[0].value, 'gemini-1.5-flash');
        assert.strictEqual(legend2.items[0].html.includes('model_1:'), true);
        assert.strictEqual(legend2.items[0].html.includes('gemini-1.5-flash'), true);
        
        assert.strictEqual(legend2.items[2].key, 'model_3');
        assert.strictEqual(legend2.items[2].value, 'claude-3-opus');
        assert.strictEqual(legend2.items[2].html.includes('model_3:'), true);
        assert.strictEqual(legend2.items[2].html.includes('claude-3-opus'), true);
        
        console.log('✅ Test Case 4.2: Hiển thị ghi chú Legend đầy đủ chính xác - ĐẠT');
    } catch (e) {
        console.error('❌ Test Case 4: Thất bại:', e.stack);
        process.exit(1);
    }

    console.log('\n🎉 TẤT CẢ CÁC BÀI KIỂM THỬ PHASE 02 ĐÃ ĐẠT 100%! 🚀');
}

runTests();
