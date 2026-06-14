# Phase 02: Frontend Multi-line UI & Parsing Logic
Status: ✅ Completed
Dependencies: [Phase 01: Backend Prefix Sanitization & Tests](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/plans/260614-1933-model-input-enhancements/phase-01-backend.md)

## Objective
Convert the single-line model input text field into a multi-line textarea that parses multiple models separated by newlines, trims whitespace, normalizes prefixes, and ignores duplicates.

## Requirements
### Functional
- Replace `<input type="text" id="new-model-input" ...>` with `<textarea id="new-model-input" ...>` in the layout.
- When the user pastes multi-line text (e.g. 5 lines of model paths), clicking "Add Model" parses all lines.
- Each line is split by newlines, stripped of whitespace, has any leading `models/` trimmed, and is filtered for empty lines.
- Only new, unique model names are added to the active model list. Duplicate model names already in the list must be skipped.
- Output log feedback to the console indicating the number of models added and skipped.

### Non-Functional
- Smooth animations when adding models to the list.
- Premium styling: ensure the textarea integrates beautifully with the existing dark/glassmorphic design.

## Implementation Steps
1. Modify [index.html](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/frontend/index.html#L52):
   - Replace `<input>` with `<textarea id="new-model-input" placeholder="..." rows="3"></textarea>`.
2. Modify [style.css](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/frontend/src/style.css#L254-L263):
   - Set `.add-model-group` to `flex-direction: column`.
   - Add styling for `.add-model-group textarea`.
   - Align `.add-model-group .btn-primary` to the bottom-right or set it to full-width.
3. Modify [main.js](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/frontend/src/main.js#L244-L259):
   - Update `handleAddModel` to split by newline `/\r?\n/`, trim, clean `models/` prefix, filter out empty inputs, deduplicate, and push unique names into `state.models`.
   - Clear the textarea after adding.

## Files to Create/Modify
- [index.html](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/frontend/index.html)
- [style.css](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/frontend/src/style.css)
- [main.js](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/frontend/src/main.js)

## Test Criteria
Create a Node.js-based unit test `tests/verify_multi_line_models.js` which:
- Mocks the DOM structure for model inputs.
- Simulates pasting a multi-line list containing mixed prefixes:
  ```
  models/gemini-2.5-flash
  gemini-2.5-flash-lite

  models/gemini-3.5-flash
  gemini-2.5-flash
  ```
- Asserts that:
  - Exactly 3 models are parsed: `gemini-2.5-flash`, `gemini-2.5-flash-lite`, and `gemini-3.5-flash`.
  - The leading `models/` prefix is correctly removed.
  - Duplicates (like the second `gemini-2.5-flash`) are skipped.
  - Empty lines are ignored.
- Execute using `node tests/verify_multi_line_models.js`.

---
Next Phase: [Phase 03: End-to-End Verification & Integration Tests](file:///home/skul9x/Desktop/Test_code/Model-Compare-main/plans/260614-1933-model-input-enhancements/phase-03-verification.md)
