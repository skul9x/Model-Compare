# Phase 03: End-to-End Verification & Integration Tests
Status: ✅ Completed
Dependencies: [Phase 02: Frontend Multi-line UI & Parsing Logic](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/plans/260614-1933-model-input-enhancements/phase-02-frontend.md)

## Objective
Verify the end-to-end integration of the backend sanitization and the frontend multi-line parsing logic under real/simulated conditions.

## Requirements
### Functional
- Verify that both formats of model names (`models/gemini-...` vs `gemini-...`) are resolved correctly in mock runs.
- Ensure mapping files (`mapping.json`) store sanitized model names for clean folder/file naming.
- Ensure the progress matrix and logging outputs display correctly when multiple models are parsed.

## Implementation Steps
1. Run all backend tests.
2. Run frontend mock parsing tests.
3. Run Dry Run batch processing simulation to ensure mapping and file outputs are correct.

## Test Criteria
Execute the following verification scripts and manual tests:

### 1. Automated Backend Tests
Run the Go unit tests to verify that no regressions occur:
```bash
go test -v ./...
```

### 2. Custom Frontend Verification
Run the new parsing verifier:
```bash
node tests/verify_multi_line_models.js
```

Run existing integration verifiers:
```bash
node tests/verify_phase3.js
```

### 3. Manual Live Run (Dry Run)
1. Launch the app in development mode:
   ```bash
   wails dev
   ```
2. In the "Danh sách Models" section:
   - Paste the following text block:
     ```
     models/gemini-2.5-flash
     models/gemini-2.5-flash-lite
     models/gemini-3-flash-preview
     models/gemini-3.1-flash-lite
     models/gemini-3.5-flash
     ```
   - Click "Add Model" (or press Enter/Ctrl+Enter if implemented).
   - Check that all 5 models are added to the list *without* the `models/` prefix.
3. In the Settings:
   - Set Gemini API Key to `DRY_RUN`.
   - Set a temporary output directory.
4. Add 2 text inputs.
5. Click **Bắt đầu tóm tắt**.
6. Verify that the batch run finishes successfully, updating the matrix and logs.
7. Open the output directory and verify:
   - The generated `mapping.json` lists the correct model names (without `models/` prefix).
   - Results are correctly written to file.
