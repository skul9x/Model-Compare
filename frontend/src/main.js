import './style.css';
import { 
    SelectDirectory, 
    LoadConfig, 
    SaveConfig, 
    SummarizeText, 
    SaveResult, 
    SaveMapping,
    StartBatchProcess,
    CancelBatchProcess
} from '../wailsjs/go/main/App';
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime';

// Prompts
const OLD_DEFAULT_PROMPT = 'Hãy tóm tắt đoạn văn bản sau đây thành một danh sách gạch đầu dòng được đánh số thứ tự 1, 2, 3... tập trung vào các ý chính quan trọng nhất.';
const NEW_DEFAULT_PROMPT = `Bạn là một chuyên gia phân tích và tóm tắt văn bản hiệu năng cao. Hãy tóm tắt đoạn văn bản được cung cấp đáp ứng các tiêu chuẩn sau:
1. Độ chính xác & Cô đọng: Trích xuất các ý chính, thông điệp quan trọng và số liệu cốt lõi một cách súc tích nhất, loại bỏ các chi tiết thừa thãi.
2. Định dạng rõ ràng: Trình bày kết quả dưới dạng danh sách gạch đầu dòng chuyên nghiệp.
3. Giữ nguyên ngôn ngữ gốc: Thực hiện tóm tắt bằng chính ngôn ngữ của văn bản đầu vào.
4. Tính trung thực: Không tự ý thêm thắt các thông tin, nhận định hoặc suy diễn ngoài văn bản gốc.`;

// Application State
const state = {
    apiKey: '',
    models: [],
    systemPrompt: '',
    outputDir: '',
    texts: ['']
};

// UI Elements
const apiKeyInput = document.getElementById('api-key-input');
const toggleApiKeyBtn = document.getElementById('toggle-api-key-btn');
const eyeIcon = document.getElementById('eye-icon');

const newModelInput = document.getElementById('new-model-input');
const addModelBtn = document.getElementById('add-model-btn');
const modelList = document.getElementById('model-list');

const systemPromptTextarea = document.getElementById('system-prompt-textarea');

const outputDirInput = document.getElementById('output-dir-input');
const browseDirBtn = document.getElementById('browse-dir-btn');

const textPartsCountBadge = document.getElementById('text-parts-count-badge');
const splitTextPreview = document.getElementById('split-text-preview');

const startSummaryBtn = document.getElementById('start-summary-btn');
const cancelSummaryBtn = document.getElementById('cancel-summary-btn');
const progressBarContainer = document.getElementById('progress-bar-container');
const progressStatusLabel = document.getElementById('progress-status-label');
const progressPercentLabel = document.getElementById('progress-percent-label');
const progressFillBar = document.getElementById('progress-fill-bar');

const statusGridContainer = document.getElementById('status-grid-container');
const statusGrid = document.getElementById('status-grid');

const consoleLogs = document.getElementById('console-logs');
const clearConsoleBtn = document.getElementById('clear-console-btn');

// Initialize App
document.addEventListener('DOMContentLoaded', () => {
    initApp();
});

async function initApp() {
    setupEventListeners();
    await loadConfig();
    renderTextInputList();
    updateTextsAndPreview();
}

// Load Configuration from Backend
async function loadConfig() {
    try {
        const cfg = await LoadConfig();
        
        state.apiKey = cfg.apiKey || '';
        state.models = cfg.models || ['gemini-1.5-flash', 'gemini-1.5-pro'];
        state.outputDir = cfg.outputDir || '';

        // Check and upgrade System Prompt
        const rawPrompt = cfg.systemPrompt ? cfg.systemPrompt.trim() : '';
        let needsSave = false;
        
        if (rawPrompt === '' || rawPrompt === OLD_DEFAULT_PROMPT) {
            state.systemPrompt = NEW_DEFAULT_PROMPT;
            needsSave = true;
        } else {
            state.systemPrompt = cfg.systemPrompt;
        }

        // Display on UI
        apiKeyInput.value = state.apiKey;
        systemPromptTextarea.value = state.systemPrompt;
        outputDirInput.value = state.outputDir;
        
        renderModels();
        
        // If config was empty or upgraded prompt, save defaults immediately
        if (!cfg.models || !cfg.systemPrompt || needsSave) {
            await saveAppConfig();
        }
        
        logConsole('Cấu hình ứng dụng đã tải thành công.', 'success');
    } catch (err) {
        logConsole(`Lỗi khi tải cấu hình: ${err}`, 'error');
    }
}


// Save Configuration to Backend
async function saveAppConfig() {
    try {
        const cfg = {
            apiKey: state.apiKey,
            models: state.models,
            systemPrompt: state.systemPrompt,
            outputDir: state.outputDir
        };
        await SaveConfig(cfg);
    } catch (err) {
        logConsole(`Không thể lưu cấu hình: ${err}`, 'error');
    }
}

// Setup all DOM Events
function setupEventListeners() {
    // Toggle API Key visibility
    toggleApiKeyBtn.addEventListener('click', () => {
        if (apiKeyInput.type === 'password') {
            apiKeyInput.type = 'text';
            eyeIcon.innerHTML = `<path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line>`;
        } else {
            apiKeyInput.type = 'password';
            eyeIcon.innerHTML = `<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle>`;
        }
    });

    // API Key input change
    apiKeyInput.addEventListener('input', () => {
        state.apiKey = apiKeyInput.value.trim();
        saveAppConfig();
    });

    // Add model
    addModelBtn.addEventListener('click', handleAddModel);
    newModelInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            handleAddModel();
        }
    });

    // System prompt change
    systemPromptTextarea.addEventListener('input', () => {
        state.systemPrompt = systemPromptTextarea.value;
        saveAppConfig();
    });

    // Browse Directory
    browseDirBtn.addEventListener('click', async () => {
        try {
            const selected = await SelectDirectory();
            if (selected) {
                state.outputDir = selected;
                outputDirInput.value = selected;
                await saveAppConfig();
                logConsole(`Đã chọn thư mục lưu kết quả: ${selected}`, 'success');
            }
        } catch (err) {
            logConsole(`Lỗi khi chọn thư mục: ${err}`, 'error');
        }
    });

    // Add text button listener
    const addTextBtn = document.getElementById('add-text-btn');
    if (addTextBtn) {
        addTextBtn.addEventListener('click', () => {
            state.texts.push('');
            renderTextInputList();
            updateTextsAndPreview();
            
            // Automatically focus the new input
            const inputs = document.querySelectorAll('.text-part-input');
            if (inputs.length > 0) {
                inputs[inputs.length - 1].focus();
            }
        });
    }

    // Start summary execution
    startSummaryBtn.addEventListener('click', startSummary);

    // Cancel summary execution
    cancelSummaryBtn.addEventListener('click', async () => {
        logConsole('Yêu cầu hủy tiến trình tóm tắt hàng loạt...', 'info');
        cancelSummaryBtn.disabled = true;
        cancelSummaryBtn.innerHTML = `
            <svg class="icon" style="animation: spin 1s linear infinite;" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="2" x2="12" y2="6"></line><line x1="12" y1="18" x2="12" y2="22"></line><line x1="4.93" y1="4.93" x2="7.76" y2="7.76"></line><line x1="16.24" y1="16.24" x2="19.07" y2="19.07"></line><line x1="2" y1="12" x2="6" y2="12"></line><line x1="18" y1="12" x2="22" y2="12"></line><line x1="4.93" y1="19.07" x2="7.76" y2="16.24"></line><line x1="16.24" y1="7.76" x2="19.07" y2="4.93"></line></svg>
            Đang hủy...
        `;
        try {
            await CancelBatchProcess();
        } catch (err) {
            logConsole(`Lỗi khi hủy tiến trình: ${err}`, 'error');
        }
    });

    // Clear logs
    clearConsoleBtn.addEventListener('click', () => {
        consoleLogs.innerHTML = '<div class="console-placeholder">[Log đã được dọn sạch]</div>';
    });
}

// Render Model list
function renderModels() {
    modelList.innerHTML = '';
    
    if (state.models.length === 0) {
        modelList.innerHTML = '<li class="no-data-text" style="padding:0.75rem; font-size:0.85rem;">Không có model nào. Hãy thêm ở trên!</li>';
        return;
    }

    state.models.forEach((model, index) => {
        const li = document.createElement('li');
        li.className = 'model-item';
        li.innerHTML = `
            <span class="model-name">${model}</span>
            <button class="btn-delete-model" data-index="${index}" title="Xóa model">
                <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><line x1="10" y1="11" x2="10" y2="17"></line><line x1="14" y1="11" x2="14" y2="17"></line></svg>
            </button>
        `;
        
        li.querySelector('.btn-delete-model').addEventListener('click', (e) => {
            const idx = parseInt(e.currentTarget.getAttribute('data-index'), 10);
            deleteModel(idx);
        });

        modelList.appendChild(li);
    });
    renderModelLegend();
}

function handleAddModel() {
    const name = newModelInput.value.trim();
    if (!name) return;
    
    if (state.models.includes(name)) {
        logConsole(`Model "${name}" đã tồn tại trong danh sách.`, 'error');
        newModelInput.value = '';
        return;
    }

    state.models.push(name);
    newModelInput.value = '';
    renderModels();
    saveAppConfig();
    logConsole(`Đã thêm model: ${name}`, 'success');
}

// Delete model
function deleteModel(index) {
    const deletedModel = state.models[index];
    state.models.splice(index, 1);
    renderModels();
    saveAppConfig();
    logConsole(`Đã xóa model: ${deletedModel}`, 'success');
}

// Render Text Input List
function renderTextInputList() {
    const container = document.getElementById('text-inputs-container');
    if (!container) return;
    
    container.innerHTML = '';
    
    if (state.texts.length === 0) {
        state.texts = [''];
    }

    state.texts.forEach((textVal, index) => {
        const item = document.createElement('div');
        item.className = 'text-input-item';
        
        // Hide delete button if only 1 item is left
        const hideDelete = state.texts.length === 1 ? 'hidden' : '';
        
        item.innerHTML = `
            <div class="text-input-header">
                <span class="text-input-label">Văn bản ${index + 1}</span>
                <div class="text-input-actions">
                    <button class="btn-copy-text" data-index="${index}" title="Sao chép văn bản">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                        </svg>
                    </button>
                    <button class="btn-delete-text ${hideDelete}" data-index="${index}" title="Xóa văn bản">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="3 6 5 6 21 6"></polyline>
                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                        </svg>
                    </button>
                </div>
            </div>
            <textarea class="text-part-input" data-index="${index}" placeholder="Nhập nội dung văn bản cần so sánh ở đây...">${textVal}</textarea>
        `;
        
        // Update state on input
        const textarea = item.querySelector('.text-part-input');
        textarea.addEventListener('input', (e) => {
            const idx = parseInt(e.target.getAttribute('data-index'), 10);
            state.texts[idx] = e.target.value;
            updateTextsAndPreview();
        });

        // Copy text button
        const copyBtn = item.querySelector('.btn-copy-text');
        copyBtn.addEventListener('click', async (e) => {
            const idx = parseInt(e.currentTarget.getAttribute('data-index'), 10);
            const content = state.texts[idx];
            
            try {
                await navigator.clipboard.writeText(content);
                
                // Show visual feedback (check icon)
                const btn = e.currentTarget;
                const originalHTML = btn.innerHTML;
                btn.innerHTML = `
                    <svg viewBox="0 0 24 24" fill="none" stroke="#0d9488" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                `;
                
                setTimeout(() => {
                    btn.innerHTML = originalHTML;
                }, 1500);
                
                logConsole(`Đã sao chép nội dung văn bản ${idx + 1} vào clipboard.`, 'success');
            } catch (err) {
                logConsole(`Không thể sao chép văn bản: ${err}`, 'error');
            }
        });

        // Delete text button
        const deleteBtn = item.querySelector('.btn-delete-text');
        deleteBtn.addEventListener('click', (e) => {
            const idx = parseInt(e.currentTarget.getAttribute('data-index'), 10);
            state.texts.splice(idx, 1);
            
            renderTextInputList();
            updateTextsAndPreview();
            logConsole(`Đã xóa văn bản thứ ${idx + 1}.`, 'info');
        });

        container.appendChild(item);
    });
}

// Update texts count and preview
function updateTextsAndPreview() {
    const count = state.texts.length;
    if (textPartsCountBadge) {
        textPartsCountBadge.textContent = `${count} văn bản`;
    }
    renderPreview();
}

// Render model legend mapping
function renderModelLegend() {
    const legendContainer = document.getElementById('model-legend-container');
    const legendList = document.getElementById('model-legend-list');
    
    if (!legendContainer || !legendList) return;
    
    legendList.innerHTML = '';
    
    if (state.models.length === 0) {
        legendContainer.classList.add('hidden');
        return;
    }
    
    legendContainer.classList.remove('hidden');
    
    state.models.forEach((model, index) => {
        const li = document.createElement('li');
        li.className = 'legend-item';
        li.innerHTML = `<span class="legend-key"><strong>model_${index + 1}:</strong></span> <span class="legend-val">${model}</span>`;
        legendList.appendChild(li);
    });
}


function renderPreview() {
    splitTextPreview.innerHTML = '';
    
    if (state.texts.length === 0) {
        splitTextPreview.innerHTML = '<div class="no-data-text">Không có văn bản hợp lệ để hiển thị preview.</div>';
        return;
    }

    state.texts.forEach((txt, idx) => {
        const previewItem = document.createElement('div');
        previewItem.className = 'preview-item';
        
        // Show first 100 characters as excerpt
        const excerpt = txt.length > 100 ? txt.substring(0, 100) + '...' : txt;
        
        previewItem.innerHTML = `
            <div class="preview-item-header">
                <span>vb_${idx + 1}</span>
                <span class="preview-item-char-count">${txt.length} ký tự</span>
            </div>
            <div class="preview-item-body">${escapeHTML(excerpt)}</div>
        `;
        splitTextPreview.appendChild(previewItem);
    });
}

// Console Logging helper
function logConsole(message, type = 'info') {
    // Remove placeholder
    const placeholder = consoleLogs.querySelector('.console-placeholder');
    if (placeholder) {
        consoleLogs.innerHTML = '';
    }

    const timestamp = new Date().toLocaleTimeString();
    const line = document.createElement('div');
    line.className = `console-line ${type}`;
    line.textContent = `[${timestamp}] ${message}`;
    
    consoleLogs.appendChild(line);
    consoleLogs.scrollTop = consoleLogs.scrollHeight;
}

// Escape HTML
function escapeHTML(str) {
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

// Start core summary loop using Wails bidirectional Events
// Start core summary loop using Wails bidirectional Events
async function startSummary() {
    // Validations
    if (!state.apiKey) {
        logConsole('Yêu cầu bắt buộc: Chưa nhập Gemini API Key!', 'error');
        alert('Vui lòng nhập API Key trước khi bắt đầu!');
        return;
    }
    if (state.models.length === 0) {
        logConsole('Yêu cầu bắt buộc: Chưa cấu hình Model nào!', 'error');
        alert('Vui lòng thêm ít nhất một Model!');
        return;
    }

    // Filter empty/whitespace-only texts
    const validTexts = state.texts.map(t => t.trim()).filter(t => t.length > 0);
    if (validTexts.length === 0) {
        logConsole('Yêu cầu bắt buộc: Không có nội dung văn bản hợp lệ để tóm tắt!', 'error');
        alert('Vui lòng nhập nội dung văn bản hợp lệ cần tóm tắt!');
        return;
    }

    if (!state.outputDir) {
        logConsole('Yêu cầu bắt buộc: Chưa chọn thư mục đầu ra!', 'error');
        alert('Vui lòng chọn thư mục lưu kết quả!');
        return;
    }

    // Prepare UI state
    startSummaryBtn.disabled = true;
    startSummaryBtn.innerHTML = `
        <svg class="icon" style="animation: spin 1s linear infinite;" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="2" x2="12" y2="6"></line><line x1="12" y1="18" x2="12" y2="22"></line><line x1="4.93" y1="4.93" x2="7.76" y2="7.76"></line><line x1="16.24" y1="16.24" x2="19.07" y2="19.07"></line><line x1="2" y1="12" x2="6" y2="12"></line><line x1="18" y1="12" x2="22" y2="12"></line><line x1="4.93" y1="19.07" x2="7.76" y2="16.24"></line><line x1="16.24" y1="7.76" x2="19.07" y2="4.93"></line></svg>
        Đang chạy tóm tắt...
    `;
    
    // Show cancel button
    cancelSummaryBtn.classList.remove('hidden');
    cancelSummaryBtn.disabled = false;
    cancelSummaryBtn.innerHTML = `
        <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect></svg>
        Hủy tiến trình
    `;
    
    progressBarContainer.classList.remove('hidden');
    statusGridContainer.classList.remove('hidden');
    
    progressFillBar.style.width = '0%';
    progressPercentLabel.textContent = '0%';
    progressStatusLabel.textContent = 'Đang chuẩn bị...';

    // Clear and build Status Grid
    statusGrid.innerHTML = '';
    
    for (let vIdx = 0; vIdx < validTexts.length; vIdx++) {
        for (let mIdx = 0; mIdx < state.models.length; mIdx++) {
            const modelName = state.models[mIdx];
            
            const taskCard = document.createElement('div');
            taskCard.className = 'task-card pending';
            taskCard.id = `task-${mIdx}-${vIdx}`;
            
            taskCard.innerHTML = `
                <div class="task-card-header">
                    <span>Văn bản ${vIdx + 1}</span>
                    <span class="task-card-status">Chờ chạy</span>
                </div>
                <div class="task-card-model" title="${modelName}">${modelName}</div>
                <div class="task-card-time">-</div>
            `;
            
            statusGrid.appendChild(taskCard);
        }
    }

    const totalTasks = validTexts.length * state.models.length;
    let completedTasks = 0;
    
    logConsole(`Bắt đầu tiến trình tóm tắt bằng Go Backend. Tổng số tác vụ: ${totalTasks} (${validTexts.length} văn bản x ${state.models.length} mô hình).`, 'info');

    // Register Wails Event Listeners
    EventsOn('task:started', (data) => {
        const { modelIdx, vbIdx } = data;
        const taskCard = document.getElementById(`task-${modelIdx}-${vbIdx}`);
        if (taskCard) {
            taskCard.className = 'task-card running';
            taskCard.querySelector('.task-card-status').textContent = 'Đang chạy';
        }
        logConsole(`Đang chạy ${state.models[modelIdx]} cho văn bản vb_${vbIdx + 1}...`, 'info');
    });

    EventsOn('task:success', (data) => {
        const { modelIdx, vbIdx, responseTimeMs } = data;
        const taskCard = document.getElementById(`task-${modelIdx}-${vbIdx}`);
        if (taskCard) {
            taskCard.className = 'task-card success';
            taskCard.querySelector('.task-card-status').textContent = 'Thành công';
            taskCard.querySelector('.task-card-time').textContent = `${(responseTimeMs / 1000).toFixed(2)}s`;
        }
        logConsole(`[Thành công] ${state.models[modelIdx]} hoàn thành vb_${vbIdx + 1} trong ${(responseTimeMs / 1000).toFixed(2)}s.`, 'success');
        
        completedTasks++;
        updateProgressBar(completedTasks, totalTasks);
    });

    EventsOn('task:failed', (data) => {
        const { modelIdx, vbIdx, errorMessage } = data;
        const taskCard = document.getElementById(`task-${modelIdx}-${vbIdx}`);
        if (taskCard) {
            taskCard.className = 'task-card error';
            taskCard.querySelector('.task-card-status').textContent = 'Thất bại';
            
            // Add error tooltip
            let tooltip = taskCard.querySelector('.tooltip-text');
            if (!tooltip) {
                tooltip = document.createElement('span');
                tooltip.className = 'tooltip-text';
                taskCard.appendChild(tooltip);
            }
            tooltip.textContent = errorMessage;
        }
        logConsole(`[Lỗi] ${state.models[modelIdx]} thất bại ở vb_${vbIdx + 1}: ${errorMessage}`, 'error');
        
        completedTasks++;
        updateProgressBar(completedTasks, totalTasks);
    });

    EventsOn('task:cancelled', () => {
        logConsole('Tiến trình tóm tắt bị hủy bởi người dùng!', 'error');
        progressStatusLabel.textContent = 'Đã hủy!';
    });

    EventsOn('task:completed', () => {
        logConsole('Tiến trình tóm tắt hàng loạt đã kết thúc.', 'success');
        if (progressStatusLabel.textContent !== 'Đã hủy!') {
            progressStatusLabel.textContent = 'Hoàn tất!';
        }
        
        // Disable cancel button & Enable start button
        cancelSummaryBtn.classList.add('hidden');
        cancelSummaryBtn.disabled = false;
        cancelSummaryBtn.innerHTML = `
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect></svg>
            Hủy tiến trình
        `;
        
        startSummaryBtn.disabled = false;
        startSummaryBtn.innerHTML = `
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
            Bắt đầu tóm tắt
        `;

        // Clean up listeners
        EventsOff('task:started');
        EventsOff('task:success');
        EventsOff('task:failed');
        EventsOff('task:cancelled');
        EventsOff('task:completed');
    });

    // Start execution
    try {
        const config = {
            apiKey: state.apiKey,
            models: state.models,
            systemPrompt: state.systemPrompt,
            outputDir: state.outputDir
        };
        await StartBatchProcess(config, validTexts);
    } catch (err) {
        logConsole(`Không thể khởi chạy tiến trình: ${err}`, 'error');
        // Reset UI if launching failed
        cancelSummaryBtn.classList.add('hidden');
        startSummaryBtn.disabled = false;
        startSummaryBtn.innerHTML = `
            <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
            Bắt đầu tóm tắt
        `;
    }
}

function updateProgressBar(completed, total) {
    const percent = Math.round((completed / total) * 100);
    progressFillBar.style.width = `${percent}%`;
    progressPercentLabel.textContent = `${percent}%`;
    progressStatusLabel.textContent = `Đang chạy (${completed}/${total} tác vụ hoàn thành)...`;
}
