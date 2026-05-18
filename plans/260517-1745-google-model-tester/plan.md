# Plan: Google Model Tester (Go + Wails v2)
Created: 2026-05-17T17:45:00+07:00
Status: 🟡 In Progress

## Overview
Ứng dụng desktop viết bằng Go + Wails v2 cho phép người dùng nhập API Key Gemini (Google AI Studio), cấu hình danh sách model tùy ý, nhập nhiều đoạn văn bản khác nhau để tóm tắt hàng loạt. 
Kết quả tóm tắt của từng model cho từng văn bản được lưu thành các file `.txt` riêng biệt với định dạng tên chuẩn hóa (`model_1-vb_1.txt`), chứa thông tin thời gian phản hồi (response time) và nội dung tóm tắt. Giao diện được thiết kế theo phong cách Light Mode hiện đại, mượt mà và trực quan.

## Tech Stack
- **Desktop Framework:** Wails v2 (Go + Vite/HTML/JS/CSS)
- **Backend:** Go (Golang) v1.26+
- **Frontend:** Vanilla HTML, Vanilla JS, Premium Vanilla CSS (Light Mode, Google Fonts "Inter", HSL Colors, Micro-animations)
- **APIs:** Google Gemini API (AI Studio SDK hoặc Direct HTTP Client qua `net/http`)
- **State Management:** Simple JS Client State + Go Config Persistence (lưu API Key và danh sách Models)

## Phases

| Phase | Name | Status | Progress |
|-------|------|--------|----------|
| 01 | Setup Environment & Wails Template | ⬜ Pending | 0% |
| 02 | Go Backend Core Logic (API & Files) | ⬜ Pending | 0% |
| 03 | Frontend UI Development (Premium Light Mode) | ⬜ Pending | 0% |
| 04 | Integration & Wails Event Communication | ⬜ Pending | 0% |
| 05 | Testing, Edge-Cases & Verification | ⬜ Pending | 0% |

## Quick Commands
- Start Phase 1: `/code phase-01`
- Check progress: `/next`
- Save context: `/save-brain`
