const assert = require('assert');

// 1. Core Logic to be tested
function parseBulkText(text) {
    if (!text.trim()) {
        return [];
    }
    // Split text by lines that contain exactly '---'
    const rawParts = text.split(/(?:^|\n)---\s*(?:\n|$)/);
    return rawParts
        .map(p => p.trim())
        .filter(p => p.length > 0);
}

function addModel(models, newModel) {
    const name = newModel.trim();
    if (!name) return models;
    if (models.includes(name)) {
        return models; // Keep unchanged
    }
    return [...models, name];
}

function deleteModel(models, index) {
    const nextModels = [...models];
    nextModels.splice(index, 1);
    return nextModels;
}

// 2. Unit Tests
function runTests() {
    console.log('🧪 Bắt đầu chạy các bài kiểm thử logic Frontend...');

    // Test Case 1: Parse Bulk Text with separators
    try {
        const text1 = 'Đoạn 1\n---\nĐoạn 2\n---\nĐoạn 3';
        const result1 = parseBulkText(text1);
        assert.strictEqual(result1.length, 3);
        assert.strictEqual(result1[0], 'Đoạn 1');
        assert.strictEqual(result1[1], 'Đoạn 2');
        assert.strictEqual(result1[2], 'Đoạn 3');
        console.log('✅ Test Case 1.1: Tách văn bản chuẩn - ĐẠT');
    } catch (e) {
        console.error('❌ Test Case 1.1: Thất bại:', e.message);
        process.exit(1);
    }

    try {
        const text2 = '   \n---\n  Đoạn 1  \n---\n  ';
        const result2 = parseBulkText(text2);
        assert.strictEqual(result2.length, 1);
        assert.strictEqual(result2[0], 'Đoạn 1');
        console.log('✅ Test Case 1.2: Lọc khoảng trắng và đoạn rỗng - ĐẠT');
    } catch (e) {
        console.error('❌ Test Case 1.2: Thất bại:', e.message);
        process.exit(1);
    }

    try {
        const text3 = 'Đoạn A\r\n---\r\nĐoạn B';
        const result3 = parseBulkText(text3);
        assert.strictEqual(result3.length, 2);
        assert.strictEqual(result3[0], 'Đoạn A');
        assert.strictEqual(result3[1], 'Đoạn B');
        console.log('✅ Test Case 1.3: Hỗ trợ CRLF line endings - ĐẠT');
    } catch (e) {
        console.error('❌ Test Case 1.3: Thất bại:', e.message);
        process.exit(1);
    }

    // Test Case 2: Model Management - Add Model
    try {
        let models = ['gemini-1.5-flash'];
        models = addModel(models, 'gemini-1.5-pro');
        assert.strictEqual(models.length, 2);
        assert.strictEqual(models[1], 'gemini-1.5-pro');
        
        // Try adding duplicate
        models = addModel(models, 'gemini-1.5-flash');
        assert.strictEqual(models.length, 2); // Should not increase
        console.log('✅ Test Case 2.1: Thêm model và chống trùng lặp - ĐẠT');
    } catch (e) {
        console.error('❌ Test Case 2.1: Thất bại:', e.message);
        process.exit(1);
    }

    try {
        let models = ['gemini-1.5-flash'];
        models = addModel(models, '   ');
        assert.strictEqual(models.length, 1); // Should ignore empty input
        console.log('✅ Test Case 2.2: Bỏ qua model trống - ĐẠT');
    } catch (e) {
        console.error('❌ Test Case 2.2: Thất bại:', e.message);
        process.exit(1);
    }

    // Test Case 3: Model Management - Delete Model
    try {
        let models = ['model-A', 'model-B', 'model-C'];
        models = deleteModel(models, 1); // Delete 'model-B'
        assert.strictEqual(models.length, 2);
        assert.strictEqual(models[0], 'model-A');
        assert.strictEqual(models[1], 'model-C');
        console.log('✅ Test Case 3: Xóa model theo index - ĐẠT');
    } catch (e) {
        console.error('❌ Test Case 3: Thất bại:', e.message);
        process.exit(1);
    }

    console.log('\n🎉 TẤT CẢ CÁC BÀI KIỂM THỬ LOGIC FRONTEND ĐÃ ĐẠT 100%! 🚀');
}

runTests();
