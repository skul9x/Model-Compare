# Phase 01: Backend Prefix Sanitization & Tests
Status: ✅ Completed
Dependencies: None

## Objective
Ensure the Go backend automatically sanitizes model names by stripping the `models/` prefix if present, preventing malformed API endpoints (e.g. `v1beta/models/models/gemini-3.1-flash-lite`).

## Requirements
### Functional
- The `SummarizeText` function in the `GeminiClient` struct must clean `modelName` by removing any leading `models/` prefix.
- The HTTP request URL should always be formatted as `.../v1beta/models/<cleaned_model_name>:generateContent`.

### Non-Functional
- The change must be backwards-compatible with model names without the prefix (e.g. `gemini-3.1-flash-lite`).

## Implementation Steps
1. Modify `backend/gemini.go` inside [SummarizeText](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/backend/gemini.go#L67-L75):
   - Add `cleanedModelName := strings.TrimPrefix(modelName, "models/")`
   - Use `cleanedModelName` when formatting the HTTP URL.
2. Create/Modify tests to verify the sanitization logic.

## Files to Create/Modify
- [gemini.go](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/backend/gemini.go) - Implement sanitization.

## Test Criteria
We will write a dedicated backend unit test in Go.
Create a new file `backend/gemini_test.go` that:
- Tests `SummarizeText` URL formatting behavior by mocking the HTTP request or verifying URL construction.
- Runs via `go test -v ./backend`.

Test cases:
- Case 1: `gemini-3.1-flash-lite` -> Output should be processed without errors.
- Case 2: `models/gemini-3.1-flash-lite` -> Output should be processed identically, stripping `models/`.

---
Next Phase: [Phase 02: Frontend Multi-line UI & Parsing Logic](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/plans/260614-1933-model-input-enhancements/phase-02-frontend.md)
