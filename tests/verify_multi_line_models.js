const assert = require('assert');

// Mock state
const state = {
    models: []
};

// Mock functions that handleAddModel calls
let renderModelsCalled = false;
let saveAppConfigCalled = false;
const logs = [];

function renderModels() {
    renderModelsCalled = true;
}

function saveAppConfig() {
    saveAppConfigCalled = true;
}

function logConsole(message, type) {
    logs.push({ message, type });
}

// Mock elements
const newModelInput = {
    value: ''
};

// The function we are testing (handleAddModel)
function handleAddModel() {
    const rawValue = newModelInput.value;
    if (!rawValue) return;

    const lines = rawValue.split(/\r?\n/);
    let addedCount = 0;
    let skippedCount = 0;
    const addedModels = [];

    for (let line of lines) {
        let name = line.trim();
        if (!name) continue;

        if (name.startsWith('models/')) {
            name = name.substring('models/'.length);
        }

        if (state.models.includes(name) || addedModels.includes(name)) {
            skippedCount++;
            continue;
        }

        state.models.push(name);
        addedModels.push(name);
        addedCount++;
    }

    newModelInput.value = '';

    if (addedCount > 0) {
        renderModels();
        saveAppConfig();
    }

    logConsole(`Đã thêm ${addedCount} model, bỏ qua ${skippedCount} model trùng lặp.`, 'success');
}

// Run Tests
function runTests() {
    console.log('🧪 Bắt đầu chạy các bài kiểm thử multi-line model parsing...');

    // Simulate pasting multi-line models
    newModelInput.value = `models/gemini-2.5-flash
gemini-2.5-flash-lite

models/gemini-3.5-flash
gemini-2.5-flash`;

    handleAddModel();

    try {
        // Assert: Exactly 3 models are parsed: gemini-2.5-flash, gemini-2.5-flash-lite, and gemini-3.5-flash
        assert.strictEqual(state.models.length, 3);
        assert.deepStrictEqual(state.models, [
            'gemini-2.5-flash',
            'gemini-2.5-flash-lite',
            'gemini-3.5-flash'
        ]);
        console.log('✅ Test Case 1: Tách dòng và trích xuất đúng 3 model - ĐẠT');

        // Assert: The leading models/ prefix is correctly removed
        assert.strictEqual(state.models.includes('models/gemini-2.5-flash'), false);
        assert.strictEqual(state.models.includes('models/gemini-3.5-flash'), false);
        console.log('✅ Test Case 2: Loại bỏ tiền tố "models/" chính xác - ĐẠT');

        // Assert: Duplicates are skipped and empty lines are ignored
        assert.strictEqual(newModelInput.value, ''); // Textarea cleared
        assert.strictEqual(renderModelsCalled, true);
        assert.strictEqual(saveAppConfigCalled, true);
        
        // Assert console log has the correct count
        const lastLog = logs[logs.length - 1];
        assert.ok(lastLog);
        assert.strictEqual(lastLog.type, 'success');
        assert.strictEqual(lastLog.message.includes('Đã thêm 3 model'), true);
        assert.strictEqual(lastLog.message.includes('bỏ qua 1 model trùng lặp'), true);
        console.log('✅ Test Case 3: Chống trùng lặp và bỏ qua dòng trống - ĐẠT');

        console.log('\n🎉 TẤT CẢ CÁC BÀI KIỂM THỬ ĐÃ ĐẠT 100%! 🚀');
    } catch (e) {
        console.error('❌ Kiểm thử thất bại:', e.stack);
        process.exit(1);
    }
}

runTests();
